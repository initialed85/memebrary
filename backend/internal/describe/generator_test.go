package describe

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/initialed85/membrary/backend/internal/store"
)

// completionResponse mimics an OpenAI-compatible server. llama.cpp strips hidden
// reasoning out of message.content, so a completion whose budget went entirely to
// thinking arrives as an empty string with finish_reason "length".
func completionResponse(content, finish string, completionTokens int) string {
	payload := map[string]any{
		"choices": []any{map[string]any{
			"message":       map[string]any{"content": content},
			"finish_reason": finish,
		}},
		"usage": map[string]any{"prompt_tokens": 815, "completion_tokens": completionTokens},
	}
	body, _ := json.Marshal(payload)
	return string(body)
}

// recordedRequest captures what the generator actually sent, because the whole
// bug lives in the request body.
type recordedRequest struct {
	MaxTokens        *int `json:"max_tokens"`
	ReasoningTokens  *int `json:"reasoning_budget_tokens"`
	ChatTemplateArgs struct {
		EnableThinking *bool `json:"enable_thinking"`
	} `json:"chat_template_kwargs"`
}

func newTestGenerator(t *testing.T, handler http.HandlerFunc, captured *[]recordedRequest, options Options) *Generator {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var request recordedRequest
		_ = json.Unmarshal(body, &request)
		*captured = append(*captured, request)
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	options.BaseURL = server.URL
	return New(options, nil, nil)
}

// testOptions mirrors the old broken shape: a small budget, reasoning suppressed.
func testOptions() Options {
	return Options{Workers: 1, MaxTokens: 600, ReasoningTokens: 0}
}

func testMeme() store.Meme {
	return store.Meme{ID: "test", MimeType: "image/png"}
}

const validAnswer = `{"description":"Two figures confront a glowing patch request.","hashtags":["devlife"],"text_tags":["ops"]}`

// The classifier is a thinking model, and reasoning_budget_tokens is the only
// knob that actually suppresses it. If the field ever stops being sent, the
// truncation failures come back, so this is asserted rather than assumed.
func TestRequestSuppressesReasoningByDefault(t *testing.T) {
	var requests []recordedRequest
	generator := newTestGenerator(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, completionResponse(validAnswer, "stop", 101))
	}, &requests, testOptions())

	if _, err := generator.requestWithEscalation(context.Background(), testMeme(), []byte("png"), "", nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if len(requests) != 1 {
		t.Fatalf("expected one request, got %d", len(requests))
	}
	if requests[0].ReasoningTokens == nil || *requests[0].ReasoningTokens != 0 {
		t.Fatalf("reasoning_budget_tokens must be sent as 0, got %v", requests[0].ReasoningTokens)
	}
}

// A server that does not understand the field must not receive it: an unknown
// body parameter is a 400 on some OpenAI-compatible providers.
func TestNegativeReasoningTokensOmitsTheField(t *testing.T) {
	var requests []recordedRequest
	options := testOptions()
	options.ReasoningTokens = -1
	generator := newTestGenerator(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, completionResponse(validAnswer, "stop", 101))
	}, &requests, options)

	if _, err := generator.requestWithEscalation(context.Background(), testMeme(), []byte("png"), "", nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if requests[0].ReasoningTokens != nil {
		t.Fatalf("reasoning_budget_tokens should be omitted, got %v", *requests[0].ReasoningTokens)
	}
}

// A model that ignores the reasoning budget and needs ~626 tokens (thinking plus
// answer) must be rescued by the escalation retry rather than failing on the
// initial budget.
func TestRequestWithEscalationRecoversFromTruncatedReasoning(t *testing.T) {
	var requests []recordedRequest
	calls := 0
	generator := newTestGenerator(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		calls++
		if calls == 1 {
			_, _ = io.WriteString(w, completionResponse("", "length", testOptions().MaxTokens))
			return
		}
		_, _ = io.WriteString(w, completionResponse(validAnswer, "stop", 626))
	}, &requests, testOptions())

	result, err := generator.requestWithEscalation(context.Background(), testMeme(), []byte("png"), "", nil)
	if err != nil {
		t.Fatalf("escalation did not recover the job: %v", err)
	}
	if result.Description == "" || len(result.Hashtags) == 0 {
		t.Fatalf("unexpected result after retry: %+v", result)
	}
	if len(requests) != 2 {
		t.Fatalf("expected exactly one retry, got %d requests", len(requests))
	}
	if *requests[0].MaxTokens != 600 || *requests[1].MaxTokens != 1200 {
		t.Fatalf("expected the budget to double once, got %d then %d", *requests[0].MaxTokens, *requests[1].MaxTokens)
	}
}

// An unrecoverable truncation has to say what happened: the old
// "unexpected end of JSON input" sent two people chasing a context-length theory
// that was not the cause.
func TestTruncatedReasoningReportsBudgetNotJSON(t *testing.T) {
	var requests []recordedRequest
	generator := newTestGenerator(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, completionResponse("", "length", 600))
	}, &requests, testOptions())

	_, err := generator.requestWithEscalation(context.Background(), testMeme(), []byte("png"), "", nil)
	var truncated *truncatedError
	if !errors.As(err, &truncated) {
		t.Fatalf("expected a truncatedError, got %v", err)
	}
	for _, want := range []string{"completion budget", "finish_reason=length", "AI_MAX_TOKENS", "AI_REASONING_TOKENS"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q should mention %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Fatalf("error should not degrade to a JSON parse message: %v", err)
	}
	if len(requests) != maxTokenEscalations+1 {
		t.Fatalf("expected the retry ceiling to be honoured, got %d requests", len(requests))
	}
}

// Non-truncation failures must not burn extra requests.
func TestOrdinaryParseFailureDoesNotEscalate(t *testing.T) {
	var requests []recordedRequest
	generator := newTestGenerator(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, completionResponse("", "stop", 1))
	}, &requests, testOptions())

	err := func() error {
		_, err := generator.requestWithEscalation(context.Background(), testMeme(), []byte("png"), "", nil)
		return err
	}()
	if err == nil {
		t.Fatal("expected an error for an empty, non-truncated completion")
	}
	if len(requests) != 1 {
		t.Fatalf("expected a single request, got %d", len(requests))
	}
	if !strings.Contains(err.Error(), "empty completion") {
		t.Fatalf("an empty non-truncated completion should say it is empty, not claim a budget problem: %v", err)
	}
	var truncated *truncatedError
	if errors.As(err, &truncated) {
		t.Fatal("an empty completion with finish_reason=stop is not a budget problem")
	}
}

// A completion cut off mid-JSON still has to yield the fields it did reach, which
// is what keeps a partial answer searchable instead of failed.
func TestPartialJSONSurvivesTruncation(t *testing.T) {
	var requests []recordedRequest
	generator := newTestGenerator(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, completionResponse(`{"description":"A wide shot of a dog.","hashtags":["dog",`, "length", 600))
	}, &requests, testOptions())

	result, err := generator.requestWithEscalation(context.Background(), testMeme(), []byte("png"), "", nil)
	if err != nil {
		t.Fatalf("partial content should be rescued, got %v", err)
	}
	if result.Description != "A wide shot of a dog." {
		t.Fatalf("unexpected description: %q", result.Description)
	}
	if len(requests) != 1 {
		t.Fatalf("a rescuable completion should not be retried, got %d requests", len(requests))
	}
}

func TestNewDefaultsMaxTokens(t *testing.T) {
	generator := New(Options{Workers: 1, MaxTokens: 0, ReasoningTokens: 0}, nil, nil)
	if generator.maxTokens != 2048 {
		t.Fatalf("expected the safe default budget, got %d", generator.maxTokens)
	}
}
