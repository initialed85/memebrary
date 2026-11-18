package describe

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/initialed85/membrary/backend/internal/store"
)

type Generator struct {
	baseURL           string
	model             string
	apiKey            string
	client            *http.Client
	store             *store.Store
	reprocessExisting bool
	workers           int
	maxTokens         int
	reasoningTokens   int
	logger            *slog.Logger
	jobs              chan store.Meme
}

// maxTokenCeiling bounds how far an escalated completion budget may grow and
// maxTokenEscalations bounds how many retries one job may cost. llama.cpp
// clamps the request to whatever context is left, so a generous ceiling is safe.
const (
	maxTokenCeiling     = 8192
	maxTokenEscalations = 2
)

// truncatedError reports that the model ran out of completion budget. A small
// thinking model returns HTTP 200 with an empty message.content and puts every
// token into hidden reasoning, which otherwise surfaces as the meaningless
// "unexpected end of JSON input".
type truncatedError struct {
	budget     int
	finish     string
	completion int
	prompt     int
}

func (e *truncatedError) Error() string {
	return fmt.Sprintf("vision model exhausted its %d-token completion budget (finish_reason=%s, completion_tokens=%d, prompt_tokens=%d); raise AI_MAX_TOKENS or set AI_REASONING_TOKENS=0 to suppress hidden reasoning",
		e.budget, e.finish, e.completion, e.prompt)
}

// Options configures one vision generator. It replaced a nine-positional-argument
// constructor, which was no longer readable.
type Options struct {
	BaseURL           string
	Model             string
	APIKey            string
	ReprocessExisting bool
	Workers           int
	MaxTokens         int
	// ReasoningTokens is the llama.cpp per-request reasoning budget. 0 stops the
	// model thinking before it answers; a negative value omits the field for a
	// server that does not understand it.
	ReasoningTokens int
}

type GeneratedContent struct {
	Description string   `json:"description"`
	Hashtags    []string `json:"hashtags"`
	TextTags    []string `json:"text_tags"`
}

func New(options Options, dataStore *store.Store, logger *slog.Logger) *Generator {
	if logger == nil {
		logger = slog.Default()
	}
	maxTokens := options.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 2048
	}
	return &Generator{
		baseURL: strings.TrimRight(options.BaseURL, "/"),
		model:   options.Model,
		apiKey:  options.APIKey,
		// A thinking completion is slow enough on a shared GPU to need more
		// headroom than the old 90s; a suppressed one answers in seconds.
		client:            &http.Client{Timeout: 180 * time.Second},
		store:             dataStore,
		reprocessExisting: options.ReprocessExisting,
		workers:           max(1, options.Workers),
		maxTokens:         maxTokens,
		reasoningTokens:   options.ReasoningTokens,
		logger:            logger,
		jobs:              make(chan store.Meme, 64),
	}
}

func (g *Generator) Enabled() bool { return g.baseURL != "" }

func (g *Generator) Start(ctx context.Context) {
	if !g.Enabled() {
		g.logger.Info("vision generator disabled")
		return
	}
	g.logger.Info("vision generator starting", "workers", g.workers, "model", g.model, "reprocess_existing", g.reprocessExisting)
	for i := 0; i < g.workers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case meme := <-g.jobs:
					g.generate(ctx, meme)
				}
			}
		}()
	}
	if pending, err := g.store.Pending(ctx); err == nil {
		g.logger.Info("queued pending vision jobs", "count", len(pending))
		for _, meme := range pending {
			g.Enqueue(meme)
		}
	} else {
		g.logger.Error("list pending vision jobs failed", "error", err)
	}
	if g.reprocessExisting {
		if allMemes, err := g.store.AllMetadata(ctx); err == nil {
			g.logger.Info("queued metadata reprocessing jobs", "count", len(allMemes))
			for _, meme := range allMemes {
				g.Enqueue(meme)
			}
		} else {
			g.logger.Error("list metadata reprocessing jobs failed", "error", err)
		}
	}
}

func (g *Generator) Enqueue(meme store.Meme) {
	if !g.Enabled() {
		return
	}
	select {
	case g.jobs <- meme:
		g.logger.Debug("vision job queued", "id", meme.ID, "tags", len(meme.Tags), "description_present", meme.Description != "")
	default:
		g.logger.Warn("vision queue full", "id", meme.ID)
		_ = g.store.UpdateGeneratedContent(context.Background(), meme.ID, meme.Description, "failed", meme.DescriptionGenerated, nil)
	}
}

func (g *Generator) generate(ctx context.Context, meme store.Meme) {
	started := time.Now()
	g.logger.Info("vision job started", "id", meme.ID, "mime", meme.MimeType, "size", meme.Size, "existing_tags", len(meme.Tags))
	imageBytes, err := osReadFile(meme.Filename)
	if err != nil {
		g.logger.Error("vision image read failed", "id", meme.ID, "error", err)
		g.markFailed(meme)
		return
	}
	existingDescription := meme.Description
	existingTags := meme.Tags
	if meme.ForceRegenerate {
		existingDescription = ""
		existingTags = nil
	}
	result, err := g.requestWithEscalation(ctx, meme, imageBytes, existingDescription, existingTags)
	if err != nil {
		g.logger.Error("vision request failed", "id", meme.ID, "error", err, "duration", time.Since(started))
		g.markFailed(meme)
		return
	}

	description := strings.TrimSpace(meme.Description)
	generatedDescription := false
	if meme.ForceRegenerate || description == "" {
		description = cleanText(result.Description)
		generatedDescription = description != ""
	}
	generatedTags := cleanTextTags(result.TextTags)
	if meme.ForceRegenerate || len(meme.Tags) == 0 {
		generatedTags = append(cleanTags(result.Hashtags, 8), generatedTags...)
	}
	// A missing description is still a failed generation. Any useful tags are
	// kept, so a retry can focus on the remaining missing field.
	status := "ready"
	if (meme.ForceRegenerate || strings.TrimSpace(meme.Description) == "") && description == "" {
		status = "failed"
	}
	if err := g.store.UpdateGeneratedContent(context.Background(), meme.ID, description, status, generatedDescription, generatedTags); err != nil {
		g.logger.Error("vision result persistence failed", "id", meme.ID, "error", err, "duration", time.Since(started))
		return
	}
	if g.reprocessExisting && len(generatedTags) > 0 {
		_ = g.store.ReplaceTags(context.Background(), meme.ID, generatedTags)
	}
	g.logger.Info("vision job finished", "id", meme.ID, "status", status, "description_generated", generatedDescription, "tags_generated", len(generatedTags), "duration", time.Since(started))
}

func (g *Generator) markFailed(meme store.Meme) {
	_ = g.store.UpdateGeneratedContent(context.Background(), meme.ID, meme.Description, "failed", meme.DescriptionGenerated, nil)
}

// osReadFile is a variable so the generator remains straightforward to test.
var osReadFile = func(name string) ([]byte, error) { return os.ReadFile(name) }

// requestWithEscalation retries with a doubled completion budget when the model
// spends everything on hidden reasoning. Without this a text-heavy meme fails on
// a coin flip: the answer is short, but the thinking that precedes it varies per
// image.
func (g *Generator) requestWithEscalation(ctx context.Context, meme store.Meme, imageBytes []byte, existingDescription string, existingTags []string) (GeneratedContent, error) {
	budget := g.maxTokens
	for attempt := 0; ; attempt++ {
		result, err := g.request(ctx, meme.MimeType, imageBytes, existingDescription, existingTags, budget)
		if err == nil {
			return result, nil
		}
		var truncated *truncatedError
		if !errors.As(err, &truncated) || attempt >= maxTokenEscalations || budget >= maxTokenCeiling {
			return GeneratedContent{}, err
		}
		budget = min(budget*2, maxTokenCeiling)
		g.logger.Warn("vision budget exhausted, retrying with a larger completion budget",
			"id", meme.ID, "attempt", attempt+1, "budget", budget, "finish_reason", truncated.finish,
			"completion_tokens", truncated.completion)
	}
}

func (g *Generator) request(ctx context.Context, mimeType string, image []byte, existingDescription string, existingTags []string, maxTokens int) (GeneratedContent, error) {
	imageURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(image)
	instruction := `Return only one valid JSON object, with exactly these fields: {"description":"one concise sentence","hashtags":["short-tag"],"text_tags":["visible-word"]}.
Describe what is visibly happening in the image. Hashtags must be 3-6 short lowercase visual or meme-context tags, without # or explanations. text_tags must contain only complete, clearly readable words or short phrases visible in the image, kept as close to the exact spelling as possible, lowercase, without punctuation. Never abbreviate, truncate, or guess a partial word; omit anything uncertain. Use an empty array only when no text is visible.
If an existing description is supplied, copy it exactly into description. If existing hashtags are supplied, copy them exactly into hashtags. Always inspect the image for text_tags.`
	payload := map[string]any{
		"model": g.model,
		"messages": []any{
			map[string]any{
				"role":    "system",
				"content": "You analyze images for a private meme library. Never mention that you are an AI and never wrap JSON in markdown.",
			},
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{"type": "text", "text": instruction + "\nExisting description: " + strings.TrimSpace(existingDescription) + "\nExisting hashtags: " + strings.Join(existingTags, ", ")},
					map[string]any{"type": "image_url", "image_url": map[string]any{"url": imageURL}},
				},
			},
		},
		// llama.cpp honours this per request (it needs b9982 or newer; older builds
		// silently ignore it, which is why maxTokens escalation still exists).
		"chat_template_kwargs": map[string]any{"enable_thinking": false},
		"response_format":      map[string]string{"type": "json_object"},
		// Use low-temperature sampling with a fixed seed: this avoids the
		// repetitive loops seen with greedy decoding while keeping retries stable.
		"temperature": 0.15,
		"top_k":       30,
		"top_p":       0.95,
		"seed":        42,
		"max_tokens":  maxTokens,
	}
	if g.reasoningTokens >= 0 {
		// A meme description does not need chain-of-thought: measured on the eight
		// memes that had been failing, suppressing reasoning cost 147 -> 135 capped
		// text tags while cutting the median completion from 859 tokens / 22.2s to
		// 164 tokens / 4.9s.
		payload["reasoning_budget_tokens"] = g.reasoningTokens
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return GeneratedContent{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return GeneratedContent{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if g.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+g.apiKey)
	}
	response, err := g.client.Do(req)
	if err != nil {
		return GeneratedContent{}, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 2*1024*1024))
	if err != nil {
		return GeneratedContent{}, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return GeneratedContent{}, fmt.Errorf("description endpoint returned %s", response.Status)
	}
	var completion struct {
		Choices []struct {
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return GeneratedContent{}, err
	}
	if len(completion.Choices) == 0 {
		return GeneratedContent{}, fmt.Errorf("description endpoint returned no choices")
	}
	choice := completion.Choices[0]
	content := cleanContent(choice.Message.Content)
	if content == "" {
		// An empty content with finish_reason "length" means the whole budget went
		// into hidden reasoning; the answer never started. Anything else that comes
		// back empty is the model's fault, not a budget to escalate.
		if choice.FinishReason == "length" {
			return GeneratedContent{}, &truncatedError{
				budget:     maxTokens,
				finish:     choice.FinishReason,
				completion: completion.Usage.CompletionTokens,
				prompt:     completion.Usage.PromptTokens,
			}
		}
		return GeneratedContent{}, fmt.Errorf("vision model returned an empty completion (finish_reason=%s, completion_tokens=%d, prompt_tokens=%d)",
			choice.FinishReason, completion.Usage.CompletionTokens, completion.Usage.PromptTokens)
	}
	return parseContent(choice.Message.Content)
}

var thinkBlock = regexp.MustCompile(`(?s)<think>.*?</think>`)

func parseContent(raw json.RawMessage) (GeneratedContent, error) {
	text := cleanContent(raw)
	var result GeneratedContent
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		// Some compatible servers ignore response_format and add a tiny preamble.
		start, end := strings.Index(text, "{"), strings.LastIndex(text, "}")
		if start >= 0 && end > start {
			if jsonErr := json.Unmarshal([]byte(text[start:end+1]), &result); jsonErr == nil {
				return normalizeResult(result), nil
			}
		}
		partial := parsePartialContent(text)
		if partial.Description != "" || len(partial.Hashtags) > 0 || len(partial.TextTags) > 0 {
			return partial, nil
		}
		if text == "" {
			return GeneratedContent{}, err
		}
		// Keep compatibility with a plain-sentence model response: the one call
		// still produces a useful description even if it misses the JSON contract.
		return GeneratedContent{Description: text}, nil
	}
	return normalizeResult(result), nil
}

func parsePartialContent(text string) GeneratedContent {
	result := GeneratedContent{}
	if description := partialString(text, "description"); description != "" {
		result.Description = cleanText(description)
	}
	result.Hashtags = partialArray(text, "hashtags")
	result.TextTags = cleanTextTags(partialArray(text, "text_tags"))
	return normalizeResult(result)
}

func partialString(text, field string) string {
	marker := `"` + field + `"`
	fieldStart := strings.Index(text, marker)
	if fieldStart < 0 {
		return ""
	}
	colon := strings.Index(text[fieldStart+len(marker):], ":")
	if colon < 0 {
		return ""
	}
	start := fieldStart + len(marker) + colon + 1
	for start < len(text) && (text[start] == ' ' || text[start] == '\t') {
		start++
	}
	if start >= len(text) || text[start] != '"' {
		return ""
	}
	start++
	for end := start; end < len(text); end++ {
		if text[end] == '"' && (end == start || text[end-1] != '\\') {
			value, err := strconv.Unquote(`"` + text[start:end] + `"`)
			if err == nil {
				return value
			}
			return ""
		}
	}
	return ""
}

func partialArray(text, field string) []string {
	marker := `"` + field + `"`
	fieldStart := strings.Index(text, marker)
	if fieldStart < 0 {
		return nil
	}
	arrayStart := strings.Index(text[fieldStart+len(marker):], "[")
	if arrayStart < 0 {
		return nil
	}
	text = text[fieldStart+len(marker)+arrayStart+1:]
	values := []string{}
	for i := 0; i < len(text); {
		if text[i] == ']' {
			break
		}
		if text[i] != '"' {
			i++
			continue
		}
		end := i + 1
		for end < len(text) {
			if text[end] == '"' && text[end-1] != '\\' {
				break
			}
			end++
		}
		if end >= len(text) {
			break
		}
		var value string
		if err := json.Unmarshal([]byte(text[i:end+1]), &value); err == nil {
			values = append(values, value)
		}
		i = end + 1
	}
	return values
}

func normalizeResult(result GeneratedContent) GeneratedContent {
	result.Description = cleanText(result.Description)
	result.Hashtags = cleanTags(result.Hashtags, 8)
	result.TextTags = cleanTextTags(result.TextTags)
	return result
}

func cleanContent(raw json.RawMessage) string {
	var text string
	if json.Unmarshal(raw, &text) != nil {
		var parts []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}
		if json.Unmarshal(raw, &parts) == nil {
			for _, part := range parts {
				if part.Text != "" {
					text += part.Text
				}
			}
		} else {
			text = string(raw)
		}
	}
	text = thinkBlock.ReplaceAllString(text, "")
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimSpace(strings.TrimPrefix(text, "```json"))
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimSpace(strings.TrimPrefix(text, "```"))
	}
	text = strings.TrimSuffix(strings.TrimSpace(text), "```")
	return strings.TrimSpace(text)
}

func cleanText(text string) string {
	text = strings.TrimSpace(strings.Trim(text, "\"'"))
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func cleanTextTags(tags []string) []string {
	meaningful := make([]string, 0, len(tags))
	for _, raw := range tags {
		words := strings.FieldsFunc(strings.ToLower(raw), func(r rune) bool { return unicode.IsSpace(r) || r == '-' || r == '_' })
		kept := make([]string, 0, len(words))
		for _, word := range words {
			word = strings.Trim(word, ".,!?;:'\"()[]{}")
			if len([]rune(word)) > 1 && !store.IsNoisyTag(word) {
				kept = append(kept, word)
			}
		}
		if len(kept) > 0 {
			meaningful = append(meaningful, strings.Join(kept, "-"))
		}
	}
	return cleanTags(meaningful, 32)
}

func cleanTags(tags []string, max int) []string {
	seen := map[string]bool{}
	result := make([]string, 0, max)
	for _, raw := range tags {
		tag := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(raw, "#")))
		tag = strings.ReplaceAll(tag, " ", "-")
		if tag == "" || seen[tag] || store.IsNoisyTag(tag) || len(result) >= max || len([]rune(tag)) > 40 {
			continue
		}
		valid := true
		for _, r := range tag {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
				valid = false
				break
			}
		}
		if valid {
			seen[tag] = true
			result = append(result, tag)
		}
	}
	return result
}
