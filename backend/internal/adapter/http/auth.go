package http

import (
	"encoding/json"
	"findopedia/internal/usecase/auth"
	"net/http"
)

type AuthHandler struct {
	uc *auth.UseCase
}

func NewAuthHandler(uc *auth.UseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Username == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password required")
		return
	}

	result, err := h.uc.Register(r.Context(), body.Username, body.Password)
	if err != nil {
		writeError(w, http.StatusConflict, "username already taken")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user":  map[string]interface{}{"id": result.User.ID, "username": result.User.Username},
		"token": result.Token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Username == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password required")
		return
	}

	result, err := h.uc.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user":  map[string]interface{}{"id": result.User.ID, "username": result.User.Username},
		"token": result.Token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromCtx(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.uc.GetMe(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}
