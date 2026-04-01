package ai

import (
	"context"
)

// Completer provides plain-text completion for note Q&A workflows.
type Completer interface {
	Complete(ctx context.Context, systemPrompt, userMessage string) (string, error)
}

// Complete sends a system prompt and user message to Gemini and returns plain assistant text.
func (c *GeminiClient) Complete(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	body := geminiGenerateReq{
		SystemInstruction: &geminiContent{Parts: []geminiPart{{Text: systemPrompt}}},
		Contents: []geminiTurn{
			{Role: "user", Parts: []geminiPart{{Text: userMessage}}},
		},
		GenerationConfig: &struct {
			MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
			Temperature     float64 `json:"temperature,omitempty"`
		}{MaxOutputTokens: 1200, Temperature: 0.4},
	}
	return c.generateText(ctx, body, "Complete")
}
