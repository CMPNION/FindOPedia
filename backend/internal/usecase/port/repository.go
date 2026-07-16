package port

import (
	"context"
	"findopedia/internal/entity"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
}

type ArticleRepository interface {
	FindByWikipediaID(ctx context.Context, wikipediaID int) (*entity.Article, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Article, error)
	Create(ctx context.Context, article *entity.Article) (*entity.Article, error)
}

type CollectionItem struct {
	Article   entity.Article
	ClaimedAt time.Time
}

type LeaderboardEntry struct {
	Username string
	Count    int
}

type QuizRepository interface {
	FindQuestionsByArticle(ctx context.Context, articleID int64) ([]entity.Question, error)
	SaveQuestions(ctx context.Context, questions []entity.Question) error
	FindAttempt(ctx context.Context, userID, articleID int64) (*entity.QuizAttempt, error)
	CreateAttempt(ctx context.Context, attempt *entity.QuizAttempt) (*entity.QuizAttempt, error)
	FindOwnership(ctx context.Context, articleID int64) (*entity.Ownership, error)
	CreateOwnership(ctx context.Context, ownership *entity.Ownership) (*entity.Ownership, error)
	FindCollectionByUserID(ctx context.Context, userID int64) ([]CollectionItem, error)
	GetLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error)
	GetLeaderboardByRarity(ctx context.Context, rarity string, limit int) ([]LeaderboardEntry, error)
	GetLastClaim(ctx context.Context, userID int64) (time.Time, error)
	UpsertLastClaim(ctx context.Context, userID int64) error
}
