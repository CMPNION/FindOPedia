package ai

import (
	"findopedia/internal/usecase/port"
	"fmt"
)

func NewQuizGenerator(provider, apiKey string) (port.QuizGenerator, error) {
	switch provider {
	case "openai":
		return NewOpenAI(apiKey), nil
	case "gemini":
		return NewGemini(apiKey), nil
	case "claude":
		return NewClaude(apiKey), nil
	default:
		return nil, fmt.Errorf("unknown ai provider: %s", provider)
	}
}
