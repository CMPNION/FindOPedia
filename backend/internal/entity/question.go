package entity

import "time"

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
)

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Question struct {
	ID            int64
	ArticleID     int64
	QuestionIndex int
	Type          QuestionType
	QuestionText  string
	Options       []Option
	CorrectAnswer string
	CreatedAt     time.Time
}
