package quiz

import (
	"context"
	"database/sql"
	"errors"
	"findopedia/internal/entity"
	"findopedia/internal/usecase/port"
	"strings"
	"time"
)

const (
	PassingThreshold = 100
	ClaimCooldown    = 4 * time.Hour
)

var (
	ErrAlreadyAttempted = errors.New("already attempted")
	ErrAlreadyOwned     = errors.New("article already owned")
	ErrNoQuestions      = errors.New("no questions found for article")
	ErrCooldown         = errors.New("claim cooldown active")
)

type UseCase struct {
	quiz     port.QuizRepository
	articles port.ArticleRepository
	users    port.UserRepository
}

func New(quiz port.QuizRepository, articles port.ArticleRepository, users port.UserRepository) *UseCase {
	return &UseCase{quiz: quiz, articles: articles, users: users}
}

type CooldownInfo struct {
	Active    bool
	NextClaim time.Time
}

func (uc *UseCase) CheckCooldown(ctx context.Context, userID int64) CooldownInfo {
	last, err := uc.quiz.GetLastClaim(ctx, userID)
	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return CooldownInfo{Active: false}
	}
	next := last.Add(ClaimCooldown)
	if time.Now().Before(next) {
		return CooldownInfo{Active: true, NextClaim: next}
	}
	return CooldownInfo{Active: false}
}

func (uc *UseCase) GetOrGenerateQuestions(ctx context.Context, slug string, userID int64, generator port.QuizGenerator) ([]entity.Question, *CooldownInfo, error) {
	article, err := uc.articles.FindBySlug(ctx, slug)
	if err != nil {
		return nil, nil, errors.New("article not found")
	}

	existing, err := uc.quiz.FindAttempt(ctx, userID, article.ID)
	if err == nil && existing != nil {
		return nil, nil, ErrAlreadyAttempted
	}

	cd := uc.CheckCooldown(ctx, userID)
	if cd.Active {
		return nil, &cd, ErrCooldown
	}

	questions, err := uc.quiz.FindQuestionsByArticle(ctx, article.ID)
	if err == nil && len(questions) > 0 {
		return questions, nil, nil
	}

	if generator == nil {
		return nil, nil, errors.New("ai generator required to create questions")
	}

	content := article.Content
	if len(content) > 50000 {
		content = content[:50000]
	}

	generated, err := generator.GenerateQuiz(ctx, port.QuizRequest{
		ArticleTitle:   article.Title,
		ArticleContent: content,
	})
	if err != nil {
		return nil, nil, err
	}

	for i := range generated {
		generated[i].ArticleID = article.ID
		generated[i].QuestionIndex = i
	}

	if err := uc.quiz.SaveQuestions(ctx, generated); err != nil {
		return nil, nil, err
	}

	saved, err := uc.quiz.FindQuestionsByArticle(ctx, article.ID)
	if err != nil {
		return nil, nil, err
	}
	return saved, nil, nil
}

type AttemptResult struct {
	Status           entity.AttemptStatus
	Score            int
	CorrectCount     int
	TotalCount       int
	OwnershipClaimed bool
	AlreadyOwner     *entity.Ownership
}

func (uc *UseCase) SubmitAttempt(ctx context.Context, slug string, userID int64, answers []SubmittedAnswer) (*AttemptResult, error) {
	article, err := uc.articles.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	questions, err := uc.quiz.FindQuestionsByArticle(ctx, article.ID)
	if err != nil || len(questions) == 0 {
		return nil, ErrNoQuestions
	}

	score, correct := GradeAnswers(questions, answers)

	status := entity.AttemptStatusFailed
	if score >= PassingThreshold {
		status = entity.AttemptStatusPassed
	}

	attempt := &entity.QuizAttempt{
		UserID:    userID,
		ArticleID: article.ID,
		Status:    status,
		Score:     score,
	}

	_, err = uc.quiz.CreateAttempt(ctx, attempt)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, ErrAlreadyAttempted
		}
		return nil, err
	}

	result := &AttemptResult{
		Status:       status,
		Score:        score,
		CorrectCount: correct,
		TotalCount:   len(questions),
	}

	if status == entity.AttemptStatusPassed {
		ownership := &entity.Ownership{
			ArticleID: article.ID,
			UserID:    userID,
		}
		_, err = uc.quiz.CreateOwnership(ctx, ownership)
		if err != nil {
			if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
				existing, _ := uc.quiz.FindOwnership(ctx, article.ID)
				result.AlreadyOwner = existing
			} else {
				return nil, err
			}
		} else {
			result.OwnershipClaimed = true
			uc.quiz.UpsertLastClaim(ctx, userID)
		}
	}

	return result, nil
}

func (uc *UseCase) GetCollection(ctx context.Context, username string) ([]port.CollectionItem, *entity.User, error) {
	user, err := uc.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}

	items, err := uc.quiz.FindCollectionByUserID(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return items, user, nil
}

type Leaderboards struct {
	Total     []port.LeaderboardEntry
	ByRarity  map[string][]port.LeaderboardEntry
}

func (uc *UseCase) GetLeaderboards(ctx context.Context) (*Leaderboards, error) {
	total, err := uc.quiz.GetLeaderboard(ctx, 10)
	if err != nil {
		return nil, err
	}

	rarities := []string{"common", "uncommon", "rare", "epic", "legendary"}
	byRarity := make(map[string][]port.LeaderboardEntry, len(rarities))
	for _, r := range rarities {
		entries, err := uc.quiz.GetLeaderboardByRarity(ctx, r, 10)
		if err != nil {
			return nil, err
		}
		byRarity[r] = entries
	}

	return &Leaderboards{Total: total, ByRarity: byRarity}, nil
}
