package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authH *AuthHandler,
	articleH *ArticleHandler,
	quizH *QuizHandler,
	parser tokenParser,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	requireAuth := RequireAuth(parser)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.With(requireAuth).Get("/me", authH.Me)
		})

		r.Route("/articles", func(r chi.Router) {
			r.Get("/random", articleH.Random)
			r.Get("/{slug}", articleH.BySlug)
			r.With(requireAuth).Post("/{slug}/questions", quizH.GetOrGenerateQuestions)
			r.With(requireAuth).Post("/{slug}/attempt", quizH.SubmitAttempt)
		})

		r.Get("/users/{username}/collection", quizH.GetCollection)
		r.Get("/leaderboard", quizH.GetLeaderboards)
		r.With(requireAuth).Get("/me/cooldown", quizH.GetCooldown)
	})

	return r
}
