package auth

import (
	"context"
	"demochat/config"
	"demochat/internal/repositories/users"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of the users.Repository for testing purposes.
type MockUserRepository struct {
	users map[string]*users.User
}

// NewMockUserRepository creates a new MockUserRepository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*users.User),
	}
}

// FindByUsername finds a user by username in the mock repository.
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*users.User, error) {
	user, ok := m.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CreateUser creates a new user in the mock repository.
func (m *MockUserRepository) CreateUser(ctx context.Context, user *users.User) error {
	if _, ok := m.users[user.Username]; ok {
		return errors.New("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

// FindByID finds a user by ID in the mock repository.
func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*users.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func TestAuthenticate_Success(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	repo.users["testuser"] = &users.User{ID: 1, Username: "testuser", Password: string(hashedPassword)}

	res, err := service.Authenticate(context.Background(), "testuser", password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.AccessToken == "" {
		t.Error("expected access token, got empty string")
	}

	if res.RefreshToken == "" {
		t.Error("expected refresh token, got empty string")
	}
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	_, err := service.Authenticate(context.Background(), "nonexistent", "password")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	repo.users["testuser"] = &users.User{ID: 1, Username: "testuser", Password: string(hashedPassword)}

	_, err := service.Authenticate(context.Background(), "testuser", "wrongpassword")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestRegister_Success(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	err := service.Register(context.Background(), "newuser", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user, ok := repo.users["newuser"]
	if !ok {
		t.Fatal("user was not created in the repository")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")); err != nil {
		t.Error("password was not hashed correctly")
	}
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	repo.users["existinguser"] = &users.User{ID: 1, Username: "existinguser", Password: "password"}

	err := service.Register(context.Background(), "existinguser", "password")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	repo.users["testuser"] = &users.User{ID: 1, Username: "testuser", Password: "password"}

	claims := jwt.MapClaims{
		"sub":      int64(1),
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, _ := token.SignedString([]byte(cfg.JWTSecret))

	res, err := service.RefreshToken(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.AccessToken == "" {
		t.Error("expected access token, got empty string")
	}

	if res.RefreshToken == "" {
		t.Error("expected refresh token, got empty string")
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	_, err := service.RefreshToken(context.Background(), "invalidtoken")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestFindByID_Success(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	repo.users["testuser"] = &users.User{ID: 1, Username: "testuser", Password: "password"}

	user, err := service.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 || user.Username != "testuser" {
		t.Errorf("unexpected user returned: %+v", user)
	}
}

func TestFindByID_UserNotFound(t *testing.T) {
	repo := NewMockUserRepository()
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := New(repo, cfg)

	_, err := service.FindByID(context.Background(), 99)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
