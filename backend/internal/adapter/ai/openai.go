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

type OpenAI struct {
	apiKey string
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{apiKey: apiKey}
}

func (o *OpenAI) GenerateQuiz(ctx context.Context, req port.QuizRequest) ([]entity.Question, error) {
	prompt := buildPrompt(req.ArticleTitle, req.ArticleContent)

	body, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai error %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("openai: no choices returned")
	}

	return parseQuestions(result.Choices[0].Message.Content)
}
