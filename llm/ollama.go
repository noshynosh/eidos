package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	ollamaURL = "http://localhost:11434"
)

type OllamaClient struct {
	Client *http.Client
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Let's make a generate function as well
func (l *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	generateReq := generateRequest{
		Model:  "llama3.2",
		Prompt: prompt,
		Stream: false,
	}

	reqBody, err := json.Marshal(generateReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaURL+"/api/generate", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (l *OllamaClient) Chat(
	ctx context.Context,
	newMessage string,
	oldMessages ...ChatMessage,
) (ChatMessage, error) {
	chatReq := chatRequest{
		Model: "llama3.2",
		Messages: append(oldMessages, ChatMessage{
			Role:    "user",
			Content: newMessage,
		}),
		Stream: false,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return ChatMessage{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaURL+"/api/chat", bytes.NewBuffer(reqBody))
	if err != nil {
		return ChatMessage{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return ChatMessage{}, err
	}
	defer resp.Body.Close()

	type respMsg struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	var rm respMsg

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatMessage{}, err
	}

	if err = json.Unmarshal(body, &rm); err != nil {
		return ChatMessage{}, err
	}

	return ChatMessage{
		Role:    "assistant",
		Content: rm.Message.Content,
	}, nil
}
