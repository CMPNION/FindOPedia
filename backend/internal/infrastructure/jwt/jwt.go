package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func IssueToken(userID int64, secret string, expiry time.Duration) (string, error) {
	claims := gojwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(expiry).Unix(),
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type Service struct {
	secret string
}

func NewService(secret string) *Service {
	return &Service{secret: secret}
}

func (s *Service) HashPassword(plain string) (string, error) { return HashPassword(plain) }
func (s *Service) CheckPassword(hash, plain string) bool     { return CheckPassword(hash, plain) }
func (s *Service) IssueToken(userID int64, expiry time.Duration) (string, error) {
	return IssueToken(userID, s.secret, expiry)
}
func (s *Service) ParseToken(tokenStr string) (int64, error) {
	return ParseToken(tokenStr, s.secret)
}

func ParseToken(tokenStr, secret string) (int64, error) {
	token, err := gojwt.Parse(tokenStr, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return 0, err
	}

	var userID int64
	if _, err := fmt.Sscanf(sub, "%d", &userID); err != nil {
		return 0, fmt.Errorf("parse user id: %w", err)
	}
	return userID, nil
}
