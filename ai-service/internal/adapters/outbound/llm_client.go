package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ChatTurn struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMClient interface {
	Complete(ctx context.Context, history []ChatTurn) (string, error)
}

type openAIClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewLLMClient() LLMClient {
	baseURL := strings.TrimRight(os.Getenv("AI_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &openAIClient{
		baseURL: baseURL,
		apiKey:  strings.TrimSpace(os.Getenv("AI_API_KEY")),
		model:   model,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func (c *openAIClient) Complete(ctx context.Context, history []ChatTurn) (string, error) {
	if c.apiKey == "" {
		return fallbackReply(history), nil
	}

	payload := map[string]interface{}{
		"model":       c.model,
		"temperature": 0.6,
		"max_tokens":  400,
		"messages":    history,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fallbackReply(history), nil
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("llm upstream: %s", resp.Status)
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 || strings.TrimSpace(parsed.Choices[0].Message.Content) == "" {
		return "", errors.New("empty llm response")
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func fallbackReply(history []ChatTurn) string {
	last := ""
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			last = strings.ToLower(history[i].Content)
			break
		}
	}

	switch {
	case strings.Contains(last, "friend"):
		return "Open Friends to send a request, accept one, or browse suggestions. Accepted friends show up on your feed."
	case strings.Contains(last, "post") || strings.Contains(last, "comment") || strings.Contains(last, "like"):
		return "Create a post from Home. You can like or comment there; the owner gets a notification automatically."
	case strings.Contains(last, "game") || strings.Contains(last, "slot") || strings.Contains(last, "spin"):
		return "Open Game to play the slot. Your player is created on first visit with balance 1000. Wins are calculated on the server."
	case strings.Contains(last, "noti"):
		return "The bell icon lists your latest notifications. They are created by the backend when someone likes, comments, or sends a friend request."
	default:
		return "I am the in-app assistant. I can help with friends, posts, notifications, and the slot game. Set AI_API_KEY on ai-service to enable full LLM answers."
	}
}
