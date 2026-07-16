package article

import (
	"context"
	"errors"
	"findopedia/internal/entity"
	"findopedia/internal/usecase/port"
	"strings"
)

var ErrNotFound = errors.New("article not found")

type UseCase struct {
	articles  port.ArticleRepository
	quiz      port.QuizRepository
	wikipedia port.WikipediaClient
}

func New(articles port.ArticleRepository, quiz port.QuizRepository, wikipedia port.WikipediaClient) *UseCase {
	return &UseCase{articles: articles, quiz: quiz, wikipedia: wikipedia}
}

type ArticleWithOwner struct {
	Article *entity.Article
	Owner   *entity.Ownership
}

func (uc *UseCase) GetRandom(ctx context.Context) (*ArticleWithOwner, error) {
	page, err := uc.wikipedia.FetchRandom(ctx)
	if err != nil {
		return nil, err
	}

	article, err := uc.articles.FindByWikipediaID(ctx, page.PageID)
	if err != nil {
		summary := page.Extract
		if len(summary) > 500 {
			summary = summary[:500]
		}
		article = &entity.Article{
			WikipediaID:   page.PageID,
			Title:         page.Title,
			Slug:          page.Slug,
			Content:       page.Extract,
			ContentLength: len(page.Extract),
			RarityTier:    entity.ComputeRarity(len(page.Extract)),
			Summary:       summary,
		}
		article, err = uc.articles.Create(ctx, article)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				article, err = uc.articles.FindByWikipediaID(ctx, page.PageID)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	owner, _ := uc.quiz.FindOwnership(ctx, article.ID)
	return &ArticleWithOwner{Article: article, Owner: owner}, nil
}

func (uc *UseCase) GetBySlug(ctx context.Context, slug string) (*ArticleWithOwner, error) {
	article, err := uc.articles.FindBySlug(ctx, slug)
	if err != nil {
		return nil, ErrNotFound
	}

	owner, _ := uc.quiz.FindOwnership(ctx, article.ID)
	return &ArticleWithOwner{Article: article, Owner: owner}, nil
}
