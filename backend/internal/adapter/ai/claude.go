package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"findopedia/internal/entity"
	"findopedia/internal/usecase/port"
	"fmt"
	"io"
	"net/http"
)

type Claude struct {
	apiKey string
}

func NewClaude(apiKey string) *Claude {
	return &Claude{apiKey: apiKey}
}

func (c *Claude) GenerateQuiz(ctx context.Context, req port.QuizRequest) ([]entity.Question, error) {
	prompt := buildPrompt(req.ArticleTitle, req.ArticleContent)

	body, _ := json.Marshal(map[string]interface{}{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 4096,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("claude error %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("claude: no content returned")
	}

	return parseQuestions(result.Content[0].Text)
}
