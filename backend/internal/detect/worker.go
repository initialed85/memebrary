package detect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/initialed85/membrary/backend/internal/store"
)

type Worker struct {
	baseURL string
	conf    float64
	client  *http.Client
	store   *store.Store
	workers int
	logger  *slog.Logger
	jobs    chan store.Meme
}

func New(baseURL string, conf float64, dataStore *store.Store, workers int, logger *slog.Logger) *Worker {
	if logger == nil {
		logger = slog.Default()
	}
	if conf <= 0 || conf > 1 {
		conf = 0.25
	}
	return &Worker{
		baseURL: strings.TrimRight(baseURL, "/"),
		conf:    conf,
		client:  &http.Client{Timeout: 60 * time.Second},
		store:   dataStore,
		workers: max(1, workers),
		logger:  logger,
		jobs:    make(chan store.Meme, 64),
	}
}

func (w *Worker) Enabled() bool { return w.baseURL != "" }

// ProcessOne runs a single job on the calling goroutine, which is what a batch
// pass needs: coral is one physical device, so the jobs must not overlap.
func (w *Worker) ProcessOne(ctx context.Context, meme store.Meme) {
	w.detect(ctx, meme)
}

func (w *Worker) Start(ctx context.Context) {
	if !w.Enabled() {
		w.logger.Info("detection worker disabled")
		return
	}
	w.logger.Info("detection worker starting", "workers", w.workers, "conf", w.conf)
	for i := 0; i < w.workers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case meme := <-w.jobs:
					w.detect(ctx, meme)
				}
			}
		}()
	}
	if pending, err := w.store.PendingDetections(ctx); err == nil {
		w.logger.Info("queued pending detection jobs", "count", len(pending))
		for _, meme := range pending {
			w.Enqueue(meme)
		}
	} else {
		w.logger.Error("list pending detection jobs failed", "error", err)
	}
}

func (w *Worker) Enqueue(meme store.Meme) {
	if !w.Enabled() {
		return
	}
	select {
	case w.jobs <- meme:
		w.logger.Debug("detection job queued", "id", meme.ID, "mime", meme.MimeType)
	default:
		w.logger.Warn("detection queue full", "id", meme.ID)
		_ = w.store.ReplaceDetections(context.Background(), meme.ID, nil, "failed")
	}
}

func (w *Worker) detect(ctx context.Context, meme store.Meme) {
	started := time.Now()
	w.logger.Info("detection job started", "id", meme.ID, "mime", meme.MimeType, "size", meme.Size)
	// Coral decodes JPEG, PNG, WebP and BMP, but not animated GIF.
	if meme.MimeType == "image/gif" {
		w.logger.Info("detection skipped for gif", "id", meme.ID)
		if err := w.store.ReplaceDetections(context.Background(), meme.ID, nil, "skipped"); err != nil {
			w.logger.Error("detection skip persistence failed", "id", meme.ID, "error", err)
		}
		return
	}
	imageBytes, err := os.ReadFile(meme.Filename)
	if err != nil {
		w.logger.Error("detection image read failed", "id", meme.ID, "error", err)
		w.markFailed(meme)
		return
	}
	result, err := w.request(ctx, meme, imageBytes)
	if err != nil {
		w.logger.Error("detection request failed", "id", meme.ID, "error", err, "duration", time.Since(started))
		w.markFailed(meme)
		return
	}
	detections := normalize(result)
	if err := w.store.ReplaceDetections(context.Background(), meme.ID, detections, "ready"); err != nil {
		w.logger.Error("detection result persistence failed", "id", meme.ID, "error", err, "duration", time.Since(started))
		return
	}
	w.logger.Info("detection job finished", "id", meme.ID, "status", "ready", "detections", len(detections), "duration", time.Since(started))
}

func (w *Worker) markFailed(meme store.Meme) {
	_ = w.store.ReplaceDetections(context.Background(), meme.ID, nil, "failed")
}

type Response struct {
	Detections []struct {
		Label      string    `json:"label"`
		ClassID    int64     `json:"class_id"`
		Confidence float64   `json:"confidence"`
		Box        []float64 `json:"box"`
	} `json:"detections"`
	ImageSize []int `json:"image_size"`
}

func (w *Worker) request(ctx context.Context, meme store.Meme, image []byte) (Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", meme.OriginalName)
	if err != nil {
		return Response{}, err
	}
	if _, err := part.Write(image); err != nil {
		return Response{}, err
	}
	if err := writer.WriteField("conf", strconv.FormatFloat(w.conf, 'f', 2, 64)); err != nil {
		return Response{}, err
	}
	if err := writer.Close(); err != nil {
		return Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.baseURL+"/detect", body)
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := w.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 2*1024*1024))
	if err != nil {
		return Response{}, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return Response{}, fmt.Errorf("detection endpoint returned %s", response.Status)
	}
	var result Response
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return Response{}, err
	}
	return result, nil
}

// normalize converts Coral's pixel-space boxes into the 0..1 space the overlay
// renders in, so a box scales with the image without the client measuring it.
func normalize(response Response) []store.Detection {
	if len(response.ImageSize) != 2 || response.ImageSize[0] <= 0 || response.ImageSize[1] <= 0 {
		return nil
	}
	width, height := float64(response.ImageSize[0]), float64(response.ImageSize[1])
	result := make([]store.Detection, 0, len(response.Detections))
	for _, detection := range response.Detections {
		if len(detection.Box) != 4 || strings.TrimSpace(detection.Label) == "" {
			continue
		}
		result = append(result, store.Detection{
			Label:      strings.TrimSpace(detection.Label),
			ClassID:    detection.ClassID,
			Confidence: detection.Confidence,
			X1:         detection.Box[0] / width,
			Y1:         detection.Box[1] / height,
			X2:         detection.Box[2] / width,
			Y2:         detection.Box[3] / height,
		})
	}
	return result
}
