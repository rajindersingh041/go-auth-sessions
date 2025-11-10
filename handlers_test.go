package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserRepo struct {
	CreateFunc func(ctx context.Context, username, passwordHash string) error
}

func (m *mockUserRepo) Create(ctx context.Context, username, passwordHash string) error {
	return m.CreateFunc(ctx, username, passwordHash)
}
func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*User, error) { return nil, nil }
func (m *mockUserRepo) FindUserID(ctx context.Context, username string) (uint64, error) { return 0, nil }
func (m *mockUserRepo) UserExists(ctx context.Context, username string) (bool, error) { return false, nil }

type mockPasswordHasher struct {
	HashFunc  func(password string) (string, error)
	CheckFunc func(password, hash string) bool
}

func (m mockPasswordHasher) Hash(password string) (string, error) { return m.HashFunc(password) }
func (m mockPasswordHasher) Check(password, hash string) bool    { return m.CheckFunc(password, hash) }

type mockJWTManager struct {
	GenerateFunc func(username string) (string, error)
	ValidateFunc func(token string) (string, error)
}

func (m mockJWTManager) Generate(username string) (string, error) { return m.GenerateFunc(username) }
func (m mockJWTManager) Validate(token string) (string, error)    { return m.ValidateFunc(token) }

func TestHandleRegister_Success(t *testing.T) {
	repo := &mockUserRepo{
		CreateFunc: func(ctx context.Context, username, passwordHash string) error {
			if username != "testuser" || passwordHash != "hashed" {
				t.Errorf("unexpected input: %s, %s", username, passwordHash)
			}
			return nil
		},
	}
	hasher := mockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			if password != "password123" {
				t.Errorf("unexpected password: %s", password)
			}
			return "hashed", nil
		},
		CheckFunc: func(password, hash string) bool { return true },
	}
	jwt := mockJWTManager{
		GenerateFunc: func(username string) (string, error) { return "token", nil },
		ValidateFunc: func(token string) (string, error) { return "testuser", nil },
	}

	srv := &Server{
		userRepository: repo,
		passwordHasher: hasher,
		jwtManager:     jwt,
		httpMux:        http.NewServeMux(),
	}
	handler := srv.handleRegister()

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
}
