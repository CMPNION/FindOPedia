package port

import (
	"context"
	"findopedia/internal/entity"
)

type QuizRequest struct {
	ArticleTitle   string
	ArticleContent string
}

type QuizGenerator interface {
	GenerateQuiz(ctx context.Context, req QuizRequest) ([]entity.Question, error)
}
