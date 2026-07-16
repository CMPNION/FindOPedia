package ai

import (
	"encoding/json"
	"findopedia/internal/entity"
	"fmt"
	"strings"
)

type rawQuestion struct {
	Index         int            `json:"index"`
	Type          string         `json:"type"`
	QuestionText  string         `json:"question_text"`
	Options       []entity.Option `json:"options"`
	CorrectAnswer string         `json:"correct_answer"`
}

func parseQuestions(raw string) ([]entity.Question, error) {
	raw = strings.TrimSpace(raw)
	// strip markdown code block if present
	if strings.HasPrefix(raw, "```") {
		lines := strings.Split(raw, "\n")
		if len(lines) > 2 {
			raw = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var items []rawQuestion
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, fmt.Errorf("parse ai response: %w", err)
	}

	questions := make([]entity.Question, 0, len(items))
	for _, item := range items {
		questions = append(questions, entity.Question{
			QuestionIndex: item.Index,
			Type:          entity.QuestionType(item.Type),
			QuestionText:  item.QuestionText,
			Options:       item.Options,
			CorrectAnswer: item.CorrectAnswer,
		})
	}
	return questions, nil
}
