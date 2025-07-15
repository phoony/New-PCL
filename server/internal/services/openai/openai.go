package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// configuration constants for OpenAI API
const (
	openAIAPIBaseURL = "https://api.openai.com/v1"
	speechModel      = "gpt-4o"
	ttsModel         = "tts-1"
	voice            = "nova"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIClient struct {
	APIKey         string
	SystemPrompt   Message
	MessageHistory []Message
	Model          string
}

// Initlaize a new OpenAI client with API key and system prompt for completions
func NewClient(apiKey string, systemPrompt string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: apiKey,
		SystemPrompt: Message{
			Role:    "system",
			Content: systemPrompt,
		},
		Model:          speechModel,
		MessageHistory: []Message{},
	}
}

// GetSingleCompletion sends only the system prompt and a single user message
func (c *OpenAIClient) GetSingleCompletion(userInput string) (string, error) {
	messages := []Message{
		c.SystemPrompt,
		{Role: "user", Content: userInput},
	}

	payload := map[string]any{
		"model":    c.Model,
		"messages": messages,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", openAIAPIBaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", body)
	}

	var result struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return result.Choices[0].Message.Content, nil
}

// TextToSpeech returns MP3 audio from input text using OpenAI TTS
func (c *OpenAIClient) TextToSpeech(text string) ([]byte, error) {
	payload := map[string]any{
		"model": ttsModel,
		"input": text,
		"voice": voice,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", openAIAPIBaseURL+"/audio/speech", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TTS error: %s", body)
	}

	return io.ReadAll(resp.Body)
}

// SetupOpenAIClient creates a new OpenAI client with the system prompt from file
func SetupOpenAIClient(apiKey string) *OpenAIClient {
	// Read system prompt from configs directory
	systemPromptPath := filepath.Join("configs", "system_prompt.md")
	systemPromptBytes, err := os.ReadFile(systemPromptPath)
	if err != nil {
		fmt.Printf("Error reading system prompt: %v\n", err)
		return nil
	}

	systemPromptContent := string(systemPromptBytes)
	return NewClient(apiKey, systemPromptContent)
}
