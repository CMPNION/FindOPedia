package main

import (
	"flag"
	"log"
	"net/http"

	adapterHTTP "findopedia/internal/adapter/http"
	"findopedia/internal/adapter/postgres"
	"findopedia/internal/adapter/wikipedia"
	"findopedia/internal/infrastructure/jwt"
	infraPostgres "findopedia/internal/infrastructure/postgres"
	"findopedia/internal/usecase/article"
	"findopedia/internal/usecase/auth"
	"findopedia/internal/usecase/quiz"
	"findopedia/shared"
)

func main() {
	envFile := flag.String("envFile", ".env", "Path to .env file")
	migrationsDir := flag.String("migrationsDir", "migrations", "Path to migrations directory")
	flag.Parse()
	if *envFile == "" {
		*envFile = ".env"
	}

	err, cfg := shared.GetConfig(*envFile)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := infraPostgres.Open(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := infraPostgres.RunMigrations(db, *migrationsDir); err != nil {
		log.Printf("migrations warning: %v", err)
	}

	jwtSvc := jwt.NewService(cfg.JWTSecret)

	userRepo := postgres.NewUserRepository(db)
	articleRepo := postgres.NewArticleRepository(db)
	quizRepo := postgres.NewQuizRepository(db)
	wikiClient := wikipedia.NewClient()

	authUC := auth.New(userRepo, jwtSvc, cfg.JWTExpiryHours)
	articleUC := article.New(articleRepo, quizRepo, wikiClient)
	quizUC := quiz.New(quizRepo, articleRepo, userRepo)

	authH := adapterHTTP.NewAuthHandler(authUC)
	articleH := adapterHTTP.NewArticleHandler(articleUC)
	quizH := adapterHTTP.NewQuizHandler(quizUC)

	router := adapterHTTP.NewRouter(authH, articleH, quizH, jwtSvc)

	addr := ":" + cfg.AppPort
	log.Printf("starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
