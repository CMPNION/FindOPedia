package auth

import (
	"context"
	"errors"
	"findopedia/internal/entity"
	"findopedia/internal/usecase/port"
	"time"
)

var (
	ErrUserExists      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type TokenIssuer interface {
	IssueToken(userID int64, expiry time.Duration) (string, error)
	CheckPassword(hash, plain string) bool
	HashPassword(plain string) (string, error)
}

type UseCase struct {
	users  port.UserRepository
	tokens TokenIssuer
	expiry time.Duration
}

func New(users port.UserRepository, tokens TokenIssuer, expiryHours int) *UseCase {
	return &UseCase{
		users:  users,
		tokens: tokens,
		expiry: time.Duration(expiryHours) * time.Hour,
	}
}

type RegisterResult struct {
	User  *entity.User
	Token string
}

func (uc *UseCase) Register(ctx context.Context, username, password string) (*RegisterResult, error) {
	hash, err := uc.tokens.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := uc.users.Create(ctx, username, hash)
	if err != nil {
		return nil, err
	}

	token, err := uc.tokens.IssueToken(user.ID, uc.expiry)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{User: user, Token: token}, nil
}

type LoginResult struct {
	User  *entity.User
	Token string
}

func (uc *UseCase) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := uc.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !uc.tokens.CheckPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	token, err := uc.tokens.IssueToken(user.ID, uc.expiry)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: user, Token: token}, nil
}

func (uc *UseCase) GetMe(ctx context.Context, userID int64) (*entity.User, error) {
	return uc.users.FindByID(ctx, userID)
}
