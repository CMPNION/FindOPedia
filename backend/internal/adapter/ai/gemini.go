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

type Gemini struct {
	apiKey string
}

func NewGemini(apiKey string) *Gemini {
	return &Gemini{apiKey: apiKey}
}

func (g *Gemini) GenerateQuiz(ctx context.Context, req port.QuizRequest) ([]entity.Question, error) {
	prompt := buildPrompt(req.ArticleTitle, req.ArticleContent)

	body, _ := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	})

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + g.apiKey
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini error %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini: no content returned")
	}

	return parseQuestions(result.Candidates[0].Content.Parts[0].Text)
}
