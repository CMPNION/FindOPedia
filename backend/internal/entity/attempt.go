package entity

import "time"

type AttemptStatus string

const (
	AttemptStatusPassed AttemptStatus = "passed"
	AttemptStatusFailed AttemptStatus = "failed"
)

type QuizAttempt struct {
	ID          int64
	UserID      int64
	ArticleID   int64
	Status      AttemptStatus
	Score       int
	AttemptedAt time.Time
}
