package ai

import "strings"

// UpstreamError captures non-2xx responses from AI providers.
type UpstreamError struct {
	Provider   string
	Operation  string
	StatusCode int
	Body       string
}

func (e *UpstreamError) Error() string {
	return e.Provider + " " + e.Operation + " upstream error"
}

// ClientMessage returns a safe, actionable message for API consumers.
func (e *UpstreamError) ClientMessage() string {
	switch e.StatusCode {
	case 400:
		return "AI provider rejected the request. Verify model configuration and request format."
	case 401, 403:
		return "AI provider authentication failed. Verify API key and provider permissions."
	case 404:
		return "AI model not found. Verify configured model id."
	case 429:
		if strings.Contains(strings.ToLower(e.Body), "quota") {
			return "AI provider quota exceeded. Check Gemini billing/quota and retry."
		}
		return "AI provider rate limit exceeded. Retry shortly."
	default:
		if e.StatusCode >= 500 {
			return "AI provider is temporarily unavailable. Retry shortly."
		}
		return "AI provider request failed."
	}
}
