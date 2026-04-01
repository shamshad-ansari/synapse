package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GeminiClient calls Google Gemini generateContent for flashcard generation.
type GeminiClient struct {
	APIKey     string
	ModelID    string
	httpClient *http.Client
}

// NewGeminiClient returns a client with default model and 60s HTTP timeout.
func NewGeminiClient(apiKey, modelID string) *GeminiClient {
	if strings.TrimSpace(modelID) == "" {
		// Use a stable "latest" alias so local dev works without explicit model pinning.
		modelID = "gemini-flash-latest"
	}
	return &GeminiClient{
		APIKey:  apiKey,
		ModelID: modelID,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type geminiGenerateReq struct {
	SystemInstruction *geminiContent `json:"systemInstruction,omitempty"`
	Contents          []geminiTurn   `json:"contents"`
	GenerationConfig  *struct {
		MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
		Temperature     float64 `json:"temperature,omitempty"`
	} `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiTurn struct {
	Role  string         `json:"role"`
	Parts []geminiPart   `json:"parts"`
}

type geminiGenerateResp struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

// GenerateFlashcards calls Gemini generateContent and parses the JSON flashcard array.
func (c *GeminiClient) GenerateFlashcards(ctx context.Context, input GenerateFlashcardsInput) ([]GeneratedCard, error) {
	system := BuildFlashcardSystemPrompt(input)
	body := geminiGenerateReq{
		SystemInstruction: &geminiContent{Parts: []geminiPart{{Text: system}}},
		Contents: []geminiTurn{
			{Role: "user", Parts: []geminiPart{{Text: input.NoteContent}}},
		},
		GenerationConfig: &struct {
			MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
			Temperature     float64 `json:"temperature,omitempty"`
		}{MaxOutputTokens: 2000, Temperature: 0.3},
	}
	text, err := c.generateText(ctx, body, "GenerateFlashcards")
	if err != nil {
		return nil, err
	}
	cards, err := ParseGeneratedCardsFromAssistantText(text)
	if err != nil {
		return nil, fmt.Errorf("gemini GenerateFlashcards: %w", err)
	}
	return cards, nil
}

func (c *GeminiClient) generateText(ctx context.Context, body geminiGenerateReq, operation string) (string, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("gemini %s: %w", operation, err)
	}
	generatePathID := strings.TrimPrefix(c.ModelID, "models/")
	u := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent",
		url.PathEscape(generatePathID),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("gemini %s: %w", operation, err)
	}
	q := req.URL.Query()
	q.Set("key", c.APIKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini %s: %w", operation, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gemini %s: %w", operation, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &UpstreamError{
			Provider:   "gemini",
			Operation:  operation,
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var parsed geminiGenerateResp
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("gemini %s: decode: %w; raw=%s", operation, err, string(respBody))
	}
	if len(parsed.Candidates) == 0 {
		return "", fmt.Errorf("gemini %s: no candidates; raw=%s", operation, string(respBody))
	}
	var text string
	for _, p := range parsed.Candidates[0].Content.Parts {
		text += p.Text
	}
	if text == "" {
		return "", fmt.Errorf("gemini %s: empty text; raw=%s", operation, string(respBody))
	}
	return text, nil
}

// GeminiEmbedClient calls Gemini embedContent with outputDimensionality matching pgvector.
type GeminiEmbedClient struct {
	APIKey     string
	ModelID    string
	httpClient *http.Client
}

// NewGeminiEmbedClient returns a client targeting gemini-embedding-001 (or custom model) with 1536-dim output.
func NewGeminiEmbedClient(apiKey, modelID string) *GeminiEmbedClient {
	if strings.TrimSpace(modelID) == "" {
		modelID = "gemini-embedding-001"
	}
	return &GeminiEmbedClient{
		APIKey:  apiKey,
		ModelID: modelID,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type geminiEmbedReq struct {
	Model                string        `json:"model"`
	Content              geminiContent `json:"content"`
	OutputDimensionality int           `json:"outputDimensionality,omitempty"`
}

type geminiEmbedResp struct {
	Embedding struct {
		Values []float64 `json:"values"`
	} `json:"embedding"`
}

// Embed returns a 1536-dimensional vector (matches VECTOR(1536) in Postgres).
func (c *GeminiEmbedClient) Embed(ctx context.Context, input EmbedInput) ([]float32, error) {
	if strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("gemini Embed: empty text")
	}
	modelName := c.ModelID
	if !strings.HasPrefix(modelName, "models/") {
		modelName = "models/" + modelName
	}
	body, err := json.Marshal(geminiEmbedReq{
		Model: modelName,
		Content: geminiContent{
			Parts: []geminiPart{{Text: input.Text}},
		},
		OutputDimensionality: ExpectedEmbeddingDims,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini Embed: %w", err)
	}
	embedPathID := strings.TrimPrefix(c.ModelID, "models/")
	u := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent",
		url.PathEscape(embedPathID),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gemini Embed: %w", err)
	}
	q := req.URL.Query()
	q.Set("key", c.APIKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini Embed: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gemini Embed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &UpstreamError{
			Provider:   "gemini",
			Operation:  "Embed",
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var parsed geminiEmbedResp
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("gemini Embed: decode: %w; raw=%s", err, string(respBody))
	}
	if len(parsed.Embedding.Values) == 0 {
		return nil, fmt.Errorf("gemini Embed: empty embedding; raw=%s", string(respBody))
	}
	if len(parsed.Embedding.Values) != ExpectedEmbeddingDims {
		return nil, fmt.Errorf("gemini Embed: expected %d dims, got %d (check model supports outputDimensionality)", ExpectedEmbeddingDims, len(parsed.Embedding.Values))
	}
	out := make([]float32, len(parsed.Embedding.Values))
	for i, v := range parsed.Embedding.Values {
		out[i] = float32(v)
	}
	return out, nil
}
