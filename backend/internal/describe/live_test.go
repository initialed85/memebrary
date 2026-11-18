package describe

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/initialed85/membrary/backend/internal/config"
	"github.com/initialed85/membrary/backend/internal/store"
)

// TestLiveClassifierCompletion is a guarded probe of the real vision endpoint. It
// is skipped unless MEMBRARY_LIVE_IMAGE is a meme file and MEMBRARY_LIVE_AI_URL is
// the OpenAI-compatible base URL, because it needs the classifier to be reachable.
//
// It exists because a model swap broke vision metadata invisibly: the endpoint kept
// returning HTTP 200, only the completion budget was wrong. Run it after any
// classifier change:
//
//	MEMBRARY_LIVE_AI_URL=https://image-classifier.dev.initialed85.cc/v1 \
//	MEMBRARY_LIVE_IMAGE=/path/to/meme.png go test ./internal/describe -run Live -v
func TestLiveClassifierCompletion(t *testing.T) {
	baseURL := os.Getenv("MEMBRARY_LIVE_AI_URL")
	imagePath := os.Getenv("MEMBRARY_LIVE_IMAGE")
	if baseURL == "" || imagePath == "" {
		t.Skip("set MEMBRARY_LIVE_AI_URL and MEMBRARY_LIVE_IMAGE to probe the real classifier")
	}
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("read probe image: %v", err)
	}

	cfg := config.Load()
	generator := New(Options{
		BaseURL:         baseURL,
		Model:           cfg.AIModel,
		APIKey:          cfg.AIAPIKey,
		Workers:         1,
		MaxTokens:       cfg.AIMaxTokens,
		ReasoningTokens: cfg.AIReasoningTokens,
	}, nil, nil)

	started := time.Now()
	result, err := generator.requestWithEscalation(context.Background(), store.Meme{ID: "probe", MimeType: "image/png"}, imageBytes, "", nil)
	if err != nil {
		t.Fatalf("live classifier request failed: %v", err)
	}
	if result.Description == "" {
		t.Fatal("live classifier returned no description")
	}
	t.Logf("probe ok in %s: %d hashtags, %d text tags, %d-char description",
		time.Since(started).Round(100*time.Millisecond), len(result.Hashtags), len(result.TextTags), len(result.Description))
}
