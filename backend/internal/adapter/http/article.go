package http

import (
	"findopedia/internal/entity"
	"findopedia/internal/usecase/article"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ArticleHandler struct {
	uc *article.UseCase
}

func NewArticleHandler(uc *article.UseCase) *ArticleHandler {
	return &ArticleHandler{uc: uc}
}

type ownerResponse struct {
	Username  string    `json:"username"`
	ClaimedAt time.Time `json:"claimed_at"`
}

type articleResponse struct {
	ID            int64             `json:"id"`
	WikipediaID   int               `json:"wikipedia_id"`
	Title         string            `json:"title"`
	Slug          string            `json:"slug"`
	Content       string            `json:"content"`
	ContentLength int               `json:"content_length"`
	RarityTier    entity.RarityTier `json:"rarity_tier"`
	Summary       string            `json:"summary"`
	Owner         *ownerResponse    `json:"owner"`
}

func toArticleResponse(art *entity.Article, owner *entity.Ownership) articleResponse {
	resp := articleResponse{
		ID:            art.ID,
		WikipediaID:   art.WikipediaID,
		Title:         art.Title,
		Slug:          art.Slug,
		Content:       art.Content,
		ContentLength: art.ContentLength,
		RarityTier:    art.RarityTier,
		Summary:       art.Summary,
	}
	if owner != nil {
		resp.Owner = &ownerResponse{
			Username:  owner.Username,
			ClaimedAt: owner.ClaimedAt,
		}
	}
	return resp
}

func (h *ArticleHandler) Random(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.GetRandom(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch article")
		return
	}
	writeJSON(w, http.StatusOK, toArticleResponse(result.Article, result.Owner))
}

func (h *ArticleHandler) BySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	result, err := h.uc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "article not found")
		return
	}
	writeJSON(w, http.StatusOK, toArticleResponse(result.Article, result.Owner))
}
