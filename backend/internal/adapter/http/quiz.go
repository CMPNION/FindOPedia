package http

import (
	"encoding/json"
	"errors"
	"findopedia/internal/adapter/ai"
	"findopedia/internal/entity"
	"findopedia/internal/usecase/port"
	"findopedia/internal/usecase/quiz"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type QuizHandler struct {
	uc *quiz.UseCase
}

func NewQuizHandler(uc *quiz.UseCase) *QuizHandler {
	return &QuizHandler{uc: uc}
}

type questionResponse struct {
	ID           int64               `json:"id"`
	Index        int                 `json:"index"`
	Type         entity.QuestionType `json:"type"`
	QuestionText string              `json:"question_text"`
	Options      []entity.Option     `json:"options"`
}

func (h *QuizHandler) GetOrGenerateQuestions(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromCtx(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	var body struct {
		AIProvider string `json:"ai_provider"`
		AIAPIKey   string `json:"ai_api_key"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	var gen port.QuizGenerator
	if body.AIProvider != "" && body.AIAPIKey != "" {
		g, err := ai.NewQuizGenerator(body.AIProvider, body.AIAPIKey)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		gen = g
	}

	questions, cd, err := h.uc.GetOrGenerateQuestions(r.Context(), slug, userID, gen)
	if err != nil {
		if errors.Is(err, quiz.ErrAlreadyAttempted) {
			writeError(w, http.StatusForbidden, "already_attempted")
			return
		}
		if errors.Is(err, quiz.ErrCooldown) && cd != nil {
			writeJSON(w, http.StatusTooManyRequests, map[string]interface{}{
				"error":      "cooldown_active",
				"next_claim": cd.NextClaim,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]questionResponse, 0, len(questions))
	for _, q := range questions {
		resp = append(resp, questionResponse{
			ID:           q.ID,
			Index:        q.QuestionIndex,
			Type:         q.Type,
			QuestionText: q.QuestionText,
			Options:      q.Options,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"questions": resp,
	})
}

func (h *QuizHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromCtx(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	var body struct {
		Answers []quiz.SubmittedAnswer `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Answers) == 0 {
		writeError(w, http.StatusBadRequest, "answers required")
		return
	}

	result, err := h.uc.SubmitAttempt(r.Context(), slug, userID, body.Answers)
	if err != nil {
		if errors.Is(err, quiz.ErrAlreadyAttempted) {
			writeError(w, http.StatusConflict, "already_attempted")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]interface{}{
		"status":            result.Status,
		"score":             result.Score,
		"correct_count":     result.CorrectCount,
		"total_count":       result.TotalCount,
		"ownership_claimed": result.OwnershipClaimed,
	}
	if result.AlreadyOwner != nil {
		resp["owner"] = map[string]interface{}{
			"username":   result.AlreadyOwner.Username,
			"claimed_at": result.AlreadyOwner.ClaimedAt,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) GetCollection(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	items, user, err := h.uc.GetCollection(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	type itemResp struct {
		Slug       string            `json:"slug"`
		Title      string            `json:"title"`
		RarityTier entity.RarityTier `json:"rarity_tier"`
		ClaimedAt  time.Time         `json:"claimed_at"`
	}

	articlesList := make([]itemResp, 0, len(items))
	for _, item := range items {
		articlesList = append(articlesList, itemResp{
			Slug:       item.Article.Slug,
			Title:      item.Article.Title,
			RarityTier: item.Article.RarityTier,
			ClaimedAt:  item.ClaimedAt,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"username": user.Username,
		"total":    len(items),
		"articles": articlesList,
	})
}

func (h *QuizHandler) GetLeaderboards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.uc.GetLeaderboards(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load leaderboard")
		return
	}
	writeJSON(w, http.StatusOK, boards)
}

func (h *QuizHandler) GetCooldown(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromCtx(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	cd := h.uc.CheckCooldown(r.Context(), userID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"active":     cd.Active,
		"next_claim": cd.NextClaim,
	})
}
