package shared

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	AppPort        string
	PostgresDSN    string
	JWTSecret      string
	JWTExpiryHours int
}

// GetConfig loads config from .env file if it exists, then reads process env vars.
// In Docker the file won't exist — env vars come from docker-compose directly.
func GetConfig(envFile string) (error, *config) {
	if envFile == "" {
		envFile = ".env"
	}
	godotenv.Load(envFile) // non-fatal: ignore missing file, fall through to process env

	jwtExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if jwtExpiry == 0 {
		jwtExpiry = 168
	}
	return nil, &config{
		AppPort:        os.Getenv("APP_PORT"),
		PostgresDSN:    os.Getenv("POSTGRES_DSN"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiryHours: jwtExpiry,
	}
}
