package detect

import (
	"encoding/json"
	"testing"

	"github.com/initialed85/membrary/backend/internal/store"
)

func TestNormalize(t *testing.T) {
	raw := `{"detections":[{"label":"person","class_id":0,"confidence":0.9115,"box":[11,21,508,597]}],"count":1,"image_size":[512,600]}`
	var response Response
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("parse coral response: %v", err)
	}
	detections := normalize(response)
	if len(detections) != 1 {
		t.Fatalf("expected one detection, got %d", len(detections))
	}
	got := detections[0]
	want := store.Detection{
		Label:      "person",
		ClassID:    0,
		Confidence: 0.9115,
		X1:         11.0 / 512.0,
		Y1:         21.0 / 600.0,
		X2:         508.0 / 512.0,
		Y2:         597.0 / 600.0,
	}
	if got != want {
		t.Fatalf("normalized box mismatch: got %+v, want %+v", got, want)
	}
}

func TestNormalizeSkipsUnusableRows(t *testing.T) {
	raw := `{"detections":[{"label":"person","class_id":0,"confidence":0.5,"box":[1,2,3]},{"label":"","class_id":1,"confidence":0.5,"box":[1,2,3,4]}],"image_size":[512,600]}`
	var response Response
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("parse coral response: %v", err)
	}
	if detections := normalize(response); len(detections) != 0 {
		t.Fatalf("expected no detections, got %d", len(detections))
	}
}

func TestNormalizeWithoutImageSize(t *testing.T) {
	raw := `{"detections":[{"label":"person","class_id":0,"confidence":0.5,"box":[1,2,3,4]}],"image_size":[]}`
	var response Response
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("parse coral response: %v", err)
	}
	if detections := normalize(response); len(detections) != 0 {
		t.Fatalf("expected no detections without image size, got %d", len(detections))
	}
}
