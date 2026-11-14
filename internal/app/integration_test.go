package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"

	"github.com/HemanthGowda/library-management-go/internal/config"
)

var application *Application

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("library_management_test"),
		postgres.WithUsername("library"),
		postgres.WithPassword("library"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres dsn: %v\n", err)
		os.Exit(1)
	}

	application, err = New(ctx, config.Config{
		DatabaseURL:       dsn,
		JWTSecret:         "test-secret-key-that-is-long-enough-for-hs256-signing",
		JWTExpiration:     time.Hour,
		JWTIssuer:         "library-management-test",
		SeedAdminEmail:    "admin@library.local",
		SeedAdminPassword: "ChangeMe123!",
		SeedAdminFullName: "Library Admin",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "application: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	application.Close()
	_ = container.Terminate(context.Background())
	os.Exit(code)
}

func TestRegisterLoginAndUpdateBookHappyPath(t *testing.T) {
	reset(t)
	insertUser(t, "admin@library.local", "ChangeMe123!", "Library Admin", "ADMIN")
	bookID := insertBook(t)

	email := "intern+" + uuid.NewString() + "@library.local"
	status, body := request(t, http.MethodPost, "/api/v1/users", "", map[string]any{
		"email":    email,
		"password": "InternPass1",
		"fullName": "New Intern",
	})
	require.Equal(t, http.StatusCreated, status)
	require.Equal(t, email, body["email"])
	roles, _ := body["roles"].([]any)
	require.Equal(t, "MEMBER", roles[0])

	token := login(t, "admin@library.local", "ChangeMe123!")
	status, body = request(t, http.MethodPut, "/api/v1/books/"+bookID, token, validBook())
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "Updated Seed Book", body["title"])
	require.Equal(t, "9780596009205", body["isbn"])
	require.EqualValues(t, 5, body["availableCopies"])
}

func TestRegisterRejectsInvalidPayload(t *testing.T) {
	reset(t)
	status, body := request(t, http.MethodPost, "/api/v1/users", "", map[string]any{
		"email":    "not-an-email",
		"password": "short",
		"fullName": "",
	})
	require.Equal(t, http.StatusBadRequest, status)
	require.Equal(t, "Validation Failed", body["error"])
	require.NotEmpty(t, body["violations"])
}

func TestUpdateBookRequiresAuthentication(t *testing.T) {
	reset(t)
	bookID := insertBook(t)
	status, _ := request(t, http.MethodPut, "/api/v1/books/"+bookID, "", validBook())
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestUpdateBookForbidsMemberRole(t *testing.T) {
	reset(t)
	insertUser(t, "member@library.local", "MemberPass1", "Library Member", "MEMBER")
	bookID := insertBook(t)
	token := login(t, "member@library.local", "MemberPass1")
	status, body := request(t, http.MethodPut, "/api/v1/books/"+bookID, token, validBook())
	require.Equal(t, http.StatusForbidden, status)
	require.Equal(t, "Forbidden", body["error"])
}

func TestRegisterAuthorForbidsDuplicateName(t *testing.T) {
	reset(t)
	insertUser(t, "admin@library.local", "ChangeMe123!", "Library Admin", "ADMIN")
	token := login(t, "admin@library.local", "ChangeMe123!")
	insertAuthor(t, "mitch albom", "I am author of book called Tuesdays with morrie", 1958, "American")
	status, _ := request(t, http.MethodPost, "/api/v1/authors", token, validAuthor())
	require.Equal(t, http.StatusConflict, status)
}

func TestLoginRejectsBadPassword(t *testing.T) {
	reset(t)
	insertUser(t, "admin@library.local", "ChangeMe123!", "Library Admin", "ADMIN")
	status, _ := request(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    "admin@library.local",
		"password": "WrongPassword1",
	})
	require.Equal(t, http.StatusUnauthorized, status)
}

func reset(t *testing.T) {
	t.Helper()
	_, err := application.Pool.Exec(context.Background(), `TRUNCATE TABLE books, authors, user_roles, users RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
}

func insertUser(t *testing.T, email, rawPassword, fullName, role string) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.MinCost)
	require.NoError(t, err)
	id := uuid.NewString()
	_, err = application.Pool.Exec(context.Background(), `
		INSERT INTO users (id, email, password_hash, full_name, enabled)
		VALUES ($1, $2, $3, $4, true)
	`, id, email, string(hash), fullName)
	require.NoError(t, err)
	_, err = application.Pool.Exec(context.Background(), `INSERT INTO user_roles (user_id, role) VALUES ($1, $2)`, id, role)
	require.NoError(t, err)
}

func insertBook(t *testing.T) string {
	t.Helper()
	id := uuid.NewString()
	_, err := application.Pool.Exec(context.Background(), `
		INSERT INTO books (id, title, author_name, isbn, description, publication_year, total_copies, available_copies)
		VALUES ($1, 'Seed Book', 'Seed Author', '9780134685991', 'A seeded book', 2018, 3, 3)
	`, id)
	require.NoError(t, err)
	return id
}

func insertAuthor(t *testing.T, name, biography string, birth int, nationality string) {
	t.Helper()
	id := uuid.NewString()
	_, err := application.Pool.Exec(context.Background(), `
		INSERT INTO authors (id, name, biography, birth, nationality)
		VALUES ($1, $2, $3, $4, $5)
	`, id, name, biography, birth, nationality)
	require.NoError(t, err)
}

func TestAdminCanAddAuthor(t *testing.T) {
	reset(t)
	insertUser(t, "admin@library.local", "ChangeMe123!", "Library Admin", "ADMIN")
	token := login(t, "admin@library.local", "ChangeMe123!")
	status, body := request(t, http.MethodPost, "/api/v1/authors", token, validAuthor())
	require.Equal(t, http.StatusCreated, status)
	require.Equal(t, "mitch albom", body["name"])
	require.Equal(t, "American", body["nationality"])
}

func TestLibrarianCanAddAuthor(t *testing.T) {
	reset(t)
	insertUser(t, "librarian@library.local", "LibrarianPass1", "Library Staff", "LIBRARIAN")
	token := login(t, "librarian@library.local", "LibrarianPass1")
	status, _ := request(t, http.MethodPost, "/api/v1/authors", token, validAuthor())
	require.Equal(t, http.StatusCreated, status)
}

func TestMemberCannotAddAuthor(t *testing.T) {
	reset(t)
	insertUser(t, "member@library.local", "MemberPass1", "Library Member", "MEMBER")
	token := login(t, "member@library.local", "MemberPass1")
	status, body := request(t, http.MethodPost, "/api/v1/authors", token, validAuthor())
	require.Equal(t, http.StatusForbidden, status)
	require.Equal(t, "Forbidden", body["error"])
}

func TestUnauthenticatedUserCannotAddAuthor(t *testing.T) {
	reset(t)
	status, _ := request(t, http.MethodPost, "/api/v1/authors", "", validAuthor())
	require.Equal(t, http.StatusUnauthorized, status)
}

func login(t *testing.T, email, rawPassword string) string {
	t.Helper()
	status, body := request(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": rawPassword,
	})
	require.Equal(t, http.StatusOK, status)
	token, _ := body["accessToken"].(string)
	require.NotEmpty(t, token)
	return token
}

func request(t *testing.T, method, path, token string, payload any) (int, map[string]any) {
	t.Helper()
	var reader *bytes.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		require.NoError(t, err)
		reader = bytes.NewReader(encoded)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, "http://localhost"+path, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := application.App.Test(req, fiber.TestConfig{Timeout: 30 * time.Second})
	require.NoError(t, err)
	defer response.Body.Close()

	var body map[string]any
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	return response.StatusCode, body
}

func validBook() map[string]any {
	return map[string]any{
		"title":           "Updated Seed Book",
		"authorName":      "Updated Author",
		"isbn":            "9780596009205",
		"description":     "Updated by integration test",
		"publicationYear": 2021,
		"totalCopies":     6,
		"availableCopies": 5,
	}
}

func validAuthor() map[string]any {
	return map[string]any{
		"name":        "mitch albom",
		"biography":   "I am author of book called Tuesdays with morrie",
		"birth":       1956,
		"nationality": "American",
	}
}
