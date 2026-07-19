package http

import (
	"golang.org/x/time/rate"

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

	// Global burst protection: 300 req/s, burst 600
	r.Use(RateLimitMiddleware(rate.NewLimiter(300, 600)))

	// Per-client: 60 req/min (1/s), burst 20; keyed by userID when authed, else by IP
	clientLimiter := NewClientLimiter(rate.Limit(1), 20)
	r.Use(clientLimiter.Middleware())

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
