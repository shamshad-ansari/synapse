package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GenerateFlashcardsInput is passed to the LLM for RAG-grounded generation.
type GenerateFlashcardsInput struct {
	NoteContent   string
	CourseContext string
	ExistingCards []CardContext
	MaxCards      int
}

// CardContext is one retrieved flashcard for deduplication context.
type CardContext struct {
	Prompt string
	Answer string
}

// GeneratedCard is one LLM-produced candidate (not yet persisted).
type GeneratedCard struct {
	Prompt   string `json:"prompt"`
	Answer   string `json:"answer"`
	CardType string `json:"card_type"`
}

// EmbedInput is text to embed.
type EmbedInput struct {
	Text string
}

// ExpectedEmbeddingDims is the pgvector column width (flashcards / note_texts).
const ExpectedEmbeddingDims = 1536

// FlashcardGenerator produces flashcard candidates from note context (Anthropic, Gemini, etc.).
type FlashcardGenerator interface {
	GenerateFlashcards(ctx context.Context, input GenerateFlashcardsInput) ([]GeneratedCard, error)
}

// TextEmbedder produces dense vectors for RAG; vectors must match ExpectedEmbeddingDims for storage.
type TextEmbedder interface {
	Embed(ctx context.Context, input EmbedInput) ([]float32, error)
}

// Client abstracts generation + embedding (implemented by pipelineClient delegating to provider-specific clients).
type Client interface {
	GenerateFlashcards(ctx context.Context, input GenerateFlashcardsInput) ([]GeneratedCard, error)
	Embed(ctx context.Context, input EmbedInput) ([]float32, error)
}

// AnthropicClient calls Claude for flashcard generation.
type AnthropicClient struct {
	APIKey     string
	ModelID    string
	httpClient *http.Client
}

// NewAnthropicClient returns a client with default model and 60s HTTP timeout.
func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		APIKey:  apiKey,
		ModelID: "claude-sonnet-4-20250514",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type anthropicMessagesReq struct {
	Model     string              `json:"model"`
	MaxTokens int                 `json:"max_tokens"`
	System    string              `json:"system"`
	Messages  []anthropicMsgBlock `json:"messages"`
}

type anthropicMsgBlock struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicMessagesResp struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// BuildFlashcardSystemPrompt builds the shared RAG system prompt for any LLM provider.
func BuildFlashcardSystemPrompt(input GenerateFlashcardsInput) string {
	maxCards := input.MaxCards
	if maxCards <= 0 {
		maxCards = 6
	}
	var existingLines strings.Builder
	for _, ec := range input.ExistingCards {
		existingLines.WriteString(fmt.Sprintf("- Q: %s | A: %s\n", ec.Prompt, ec.Answer))
	}
	existingBlock := existingLines.String()
	if existingBlock == "" {
		existingBlock = "(none)\n"
	}
	return fmt.Sprintf(`You are an expert flashcard creator for university students.
Generate exactly %d high-quality question-answer flashcard pairs
from the provided note content. Each card should test a single, specific concept.
Course: %s.
IMPORTANT: Do NOT duplicate or closely paraphrase any of the existing cards
listed below. Focus on concepts not yet covered.
Existing cards to avoid duplicating:
%s
Respond ONLY with a JSON array. No preamble, no markdown fences, no explanation.
Format: [{"prompt":"...","answer":"...","card_type":"qa"}, ...]`, maxCards, input.CourseContext, existingBlock)
}

// ParseGeneratedCardsFromAssistantText strips fences and unmarshals the JSON array.
func ParseGeneratedCardsFromAssistantText(text string) ([]GeneratedCard, error) {
	text = stripMarkdownFences(text)
	var cards []GeneratedCard
	if err := json.Unmarshal([]byte(text), &cards); err != nil {
		return nil, fmt.Errorf("parse flashcard JSON: %w; raw=%s", err, text)
	}
	return cards, nil
}

// GenerateFlashcards builds the system prompt, calls Anthropic Messages API, parses JSON array from assistant text.
func (c *AnthropicClient) GenerateFlashcards(ctx context.Context, input GenerateFlashcardsInput) ([]GeneratedCard, error) {
	system := BuildFlashcardSystemPrompt(input)

	body := anthropicMessagesReq{
		Model:     c.ModelID,
		MaxTokens: 2000,
		System:    system,
		Messages: []anthropicMsgBlock{
			{Role: "user", Content: input.NoteContent},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: %w", err)
	}
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &UpstreamError{
			Provider:   "anthropic",
			Operation:  "GenerateFlashcards",
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var parsed anthropicMessagesResp
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: decode response: %w; raw=%s", err, string(respBody))
	}
	var text string
	for _, block := range parsed.Content {
		if block.Type == "text" && block.Text != "" {
			text += block.Text
		}
	}
	if text == "" {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: empty text content; raw=%s", string(respBody))
	}

	cards, err := ParseGeneratedCardsFromAssistantText(text)
	if err != nil {
		return nil, fmt.Errorf("anthropic GenerateFlashcards: %w", err)
	}
	return cards, nil
}

// Complete sends a system prompt and user message; returns plain assistant text.
func (c *AnthropicClient) Complete(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	body := anthropicMessagesReq{
		Model:     c.ModelID,
		MaxTokens: 2048,
		System:    systemPrompt,
		Messages: []anthropicMsgBlock{
			{Role: "user", Content: userMessage},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("anthropic Complete: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("anthropic Complete: %w", err)
	}
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic Complete: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("anthropic Complete: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &UpstreamError{
			Provider:   "anthropic",
			Operation:  "Complete",
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var parsed anthropicMessagesResp
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("anthropic Complete: decode response: %w; raw=%s", err, string(respBody))
	}
	var text string
	for _, block := range parsed.Content {
		if block.Type == "text" && block.Text != "" {
			text += block.Text
		}
	}
	if text == "" {
		return "", fmt.Errorf("anthropic Complete: empty text content; raw=%s", string(respBody))
	}
	return text, nil
}

func stripMarkdownFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```JSON")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return s
}

// OpenAIEmbedClient calls OpenAI embeddings API.
type OpenAIEmbedClient struct {
	APIKey     string
	httpClient *http.Client
}

// NewOpenAIEmbedClient returns a client with 60s HTTP timeout.
func NewOpenAIEmbedClient(apiKey string) *OpenAIEmbedClient {
	return &OpenAIEmbedClient{
		APIKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type openAIEmbedReq struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type openAIEmbedResp struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

// Embed returns a 1536-dimensional vector for text-embedding-3-small.
func (c *OpenAIEmbedClient) Embed(ctx context.Context, input EmbedInput) ([]float32, error) {
	if strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("openai Embed: empty text")
	}
	body, err := json.Marshal(openAIEmbedReq{
		Model: "text-embedding-3-small",
		Input: input.Text,
	})
	if err != nil {
		return nil, fmt.Errorf("openai Embed: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai Embed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai Embed: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai Embed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &UpstreamError{
			Provider:   "openai",
			Operation:  "Embed",
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var parsed openAIEmbedResp
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("openai Embed: decode: %w; raw=%s", err, string(respBody))
	}
	if len(parsed.Data) == 0 || len(parsed.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("openai Embed: empty embedding; raw=%s", string(respBody))
	}
	out := make([]float32, len(parsed.Data[0].Embedding))
	for i, v := range parsed.Data[0].Embedding {
		out[i] = float32(v)
	}
	if len(out) != ExpectedEmbeddingDims {
		return nil, fmt.Errorf("openai Embed: expected %d dims, got %d", ExpectedEmbeddingDims, len(out))
	}
	return out, nil
}

// pipelineClient implements Client by delegating to Anthropic + OpenAI clients.
type pipelineClient struct {
	gen FlashcardGenerator
	emb TextEmbedder
}

// NewPipelineClient returns a Client that uses gen for LLM and emb for embeddings.
func NewPipelineClient(gen FlashcardGenerator, emb TextEmbedder) Client {
	return &pipelineClient{gen: gen, emb: emb}
}

func (p *pipelineClient) GenerateFlashcards(ctx context.Context, input GenerateFlashcardsInput) ([]GeneratedCard, error) {
	return p.gen.GenerateFlashcards(ctx, input)
}

func (p *pipelineClient) Embed(ctx context.Context, input EmbedInput) ([]float32, error) {
	return p.emb.Embed(ctx, input)
}