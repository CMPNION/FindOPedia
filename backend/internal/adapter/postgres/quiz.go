package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"findopedia/internal/entity"
	"findopedia/internal/usecase/port"
	"time"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) FindQuestionsByArticle(ctx context.Context, articleID int64) ([]entity.Question, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, article_id, question_index, question_type, question_text, options, correct_answer, created_at
		 FROM quiz_questions WHERE article_id = $1 ORDER BY question_index`,
		articleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var q entity.Question
		var optionsJSON []byte
		if err := rows.Scan(&q.ID, &q.ArticleID, &q.QuestionIndex, &q.Type, &q.QuestionText, &optionsJSON, &q.CorrectAnswer, &q.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *QuizRepository) SaveQuestions(ctx context.Context, questions []entity.Question) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO quiz_questions (article_id, question_index, question_type, question_text, options, correct_answer)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range questions {
		optionsJSON, err := json.Marshal(q.Options)
		if err != nil {
			return err
		}
		if _, err := stmt.ExecContext(ctx, q.ArticleID, q.QuestionIndex, q.Type, q.QuestionText, optionsJSON, q.CorrectAnswer); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *QuizRepository) FindAttempt(ctx context.Context, userID, articleID int64) (*entity.QuizAttempt, error) {
	var a entity.QuizAttempt
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, article_id, status, score, attempted_at
		 FROM quiz_attempts WHERE user_id = $1 AND article_id = $2`,
		userID, articleID,
	).Scan(&a.ID, &a.UserID, &a.ArticleID, &a.Status, &a.Score, &a.AttemptedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *QuizRepository) CreateAttempt(ctx context.Context, attempt *entity.QuizAttempt) (*entity.QuizAttempt, error) {
	var a entity.QuizAttempt
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO quiz_attempts (user_id, article_id, status, score)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, article_id, status, score, attempted_at`,
		attempt.UserID, attempt.ArticleID, attempt.Status, attempt.Score,
	).Scan(&a.ID, &a.UserID, &a.ArticleID, &a.Status, &a.Score, &a.AttemptedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *QuizRepository) FindOwnership(ctx context.Context, articleID int64) (*entity.Ownership, error) {
	var o entity.Ownership
	err := r.db.QueryRowContext(ctx,
		`SELECT o.id, o.article_id, o.user_id, u.username, o.claimed_at
		 FROM ownership o JOIN users u ON u.id = o.user_id
		 WHERE o.article_id = $1`,
		articleID,
	).Scan(&o.ID, &o.ArticleID, &o.UserID, &o.Username, &o.ClaimedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *QuizRepository) CreateOwnership(ctx context.Context, ownership *entity.Ownership) (*entity.Ownership, error) {
	var o entity.Ownership
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO ownership (article_id, user_id)
		 VALUES ($1, $2)
		 RETURNING id, article_id, user_id, claimed_at`,
		ownership.ArticleID, ownership.UserID,
	).Scan(&o.ID, &o.ArticleID, &o.UserID, &o.ClaimedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *QuizRepository) FindCollectionByUserID(ctx context.Context, userID int64) ([]port.CollectionItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT a.id, a.wikipedia_id, a.title, a.slug, a.content, a.content_length, a.rarity_tier, a.summary, a.created_at, o.claimed_at
		 FROM ownership o
		 JOIN articles a ON a.id = o.article_id
		 WHERE o.user_id = $1
		 ORDER BY o.claimed_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []port.CollectionItem
	for rows.Next() {
		var item port.CollectionItem
		if err := rows.Scan(
			&item.Article.ID, &item.Article.WikipediaID, &item.Article.Title,
			&item.Article.Slug, &item.Article.Content, &item.Article.ContentLength,
			&item.Article.RarityTier, &item.Article.Summary, &item.Article.CreatedAt,
			&item.ClaimedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *QuizRepository) GetLeaderboard(ctx context.Context, limit int) ([]port.LeaderboardEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT u.username, COUNT(o.id) AS cnt
		 FROM ownership o JOIN users u ON u.id = o.user_id
		 GROUP BY u.username ORDER BY cnt DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []port.LeaderboardEntry
	for rows.Next() {
		var e port.LeaderboardEntry
		if err := rows.Scan(&e.Username, &e.Count); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *QuizRepository) GetLeaderboardByRarity(ctx context.Context, rarity string, limit int) ([]port.LeaderboardEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT u.username, COUNT(o.id) AS cnt
		 FROM ownership o
		 JOIN users u ON u.id = o.user_id
		 JOIN articles a ON a.id = o.article_id
		 WHERE a.rarity_tier = $1
		 GROUP BY u.username ORDER BY cnt DESC LIMIT $2`,
		rarity, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []port.LeaderboardEntry
	for rows.Next() {
		var e port.LeaderboardEntry
		if err := rows.Scan(&e.Username, &e.Count); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *QuizRepository) GetLastClaim(ctx context.Context, userID int64) (time.Time, error) {
	var t time.Time
	err := r.db.QueryRowContext(ctx,
		`SELECT last_claim FROM claim_cooldowns WHERE user_id = $1`, userID,
	).Scan(&t)
	return t, err
}

func (r *QuizRepository) UpsertLastClaim(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO claim_cooldowns (user_id, last_claim) VALUES ($1, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET last_claim = NOW()`,
		userID,
	)
	return err
}
