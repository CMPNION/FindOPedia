package postgres

import (
	"context"
	"database/sql"
	"findopedia/internal/entity"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) FindByWikipediaID(ctx context.Context, wikipediaID int) (*entity.Article, error) {
	var a entity.Article
	err := r.db.QueryRowContext(ctx,
		`SELECT id, wikipedia_id, title, slug, content, content_length, rarity_tier, summary, created_at
		 FROM articles WHERE wikipedia_id = $1`,
		wikipediaID,
	).Scan(&a.ID, &a.WikipediaID, &a.Title, &a.Slug, &a.Content, &a.ContentLength, &a.RarityTier, &a.Summary, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) FindBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	var a entity.Article
	err := r.db.QueryRowContext(ctx,
		`SELECT id, wikipedia_id, title, slug, content, content_length, rarity_tier, summary, created_at
		 FROM articles WHERE slug = $1`,
		slug,
	).Scan(&a.ID, &a.WikipediaID, &a.Title, &a.Slug, &a.Content, &a.ContentLength, &a.RarityTier, &a.Summary, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	var a entity.Article
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO articles (wikipedia_id, title, slug, content, content_length, rarity_tier, summary)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, wikipedia_id, title, slug, content, content_length, rarity_tier, summary, created_at`,
		article.WikipediaID, article.Title, article.Slug, article.Content,
		article.ContentLength, article.RarityTier, article.Summary,
	).Scan(&a.ID, &a.WikipediaID, &a.Title, &a.Slug, &a.Content, &a.ContentLength, &a.RarityTier, &a.Summary, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
