package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database and applies migrations.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&Post{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

// performRequest executes an HTTP request against the provided handler.
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		req = httptest.NewRequest(method, path, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetPostsEmpty(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	w := performRequest(router, "GET", "/posts", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var posts []Post
	if err := json.Unmarshal(w.Body.Bytes(), &posts); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if len(posts) != 0 {
		t.Fatalf("expected no posts, got %d", len(posts))
	}
}

func TestCreateAndGetPost(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	input := Post{Title: "Hello", Content: "World"}
	// Create a new post
	w := performRequest(router, "POST", "/posts", input)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on create, got %d", w.Code)
	}
	var created Post
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid JSON on create: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected non-zero ID after creation")
	}
	if created.Title != input.Title || created.Content != input.Content {
		t.Fatalf("mismatched data: got %+v, want %+v", created, input)
	}
	// Retrieve the created post by ID
	path := fmt.Sprintf("/posts/%d", created.ID)
	w = performRequest(router, "GET", path, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on get, got %d", w.Code)
	}
	var fetched Post
	if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("invalid JSON on get: %v", err)
	}
	if fetched.ID != created.ID || fetched.Title != input.Title {
		t.Fatalf("fetched post mismatch: got %+v, want %+v", fetched, created)
	}
}
