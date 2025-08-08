package auth

import (
	"context"
	"demochat/config"
	"demochat/internal/repositories/users"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationResult struct {
	AccessToken  string
	RefreshToken string
}

type Service interface {
	Authenticate(ctx context.Context, username, password string) (*AuthenticationResult, error)
	Register(ctx context.Context, username, password string) error
	RefreshToken(ctx context.Context, token string) (*AuthenticationResult, error)
	FindByID(ctx context.Context, id int64) (*users.User, error)
}

type service struct {
	repo      users.Repository
	jwtSecret string
}

func New(repo users.Repository, cfg *config.Config) Service {
	return &service{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (*AuthenticationResult, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthenticationResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Register(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &users.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *service) RefreshToken(ctx context.Context, token string) (*AuthenticationResult, error) {
	claims := &jwt.MapClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, err
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		return nil, err
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthenticationResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) FindByID(ctx context.Context, id int64) (*users.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) generateAccessToken(user *users.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *service) generateRefreshToken(user *users.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}
