package main

import (
   "encoding/json"
   "net/http"
   "net/http/httptest"
   "strings"
   "testing"
)

// TestGetPostNotFound ensures requesting a nonexistent post returns 404
func TestGetPostNotFound(t *testing.T) {
   db := setupTestDB(t)
   router := setupRouter(db)
   w := performRequest(router, "GET", "/posts/999", nil)
   if w.Code != http.StatusNotFound {
       t.Fatalf("expected status 404 for nonexistent post, got %d", w.Code)
   }
   var resp map[string]string
   if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
       t.Fatalf("invalid JSON response: %v", err)
   }
   if resp["error"] != "post not found" {
       t.Fatalf("expected error 'post not found', got '%s'", resp["error"])
   }
}

// TestCreatePostBadRequest ensures invalid JSON yields a 400 error
func TestCreatePostBadRequest(t *testing.T) {
   db := setupTestDB(t)
   router := setupRouter(db)
   req := httptest.NewRequest("POST", "/posts", strings.NewReader("{invalid_json"))
   req.Header.Set("Content-Type", "application/json")
   w := httptest.NewRecorder()
   router.ServeHTTP(w, req)
   if w.Code != http.StatusBadRequest {
       t.Fatalf("expected status 400 for bad request, got %d", w.Code)
   }
}

// TestGetPostsAfterCreation ensures posts are listed in descending creation order
func TestGetPostsAfterCreation(t *testing.T) {
   db := setupTestDB(t)
   router := setupRouter(db)
   // Create two posts
   performRequest(router, "POST", "/posts", Post{Title: "First", Content: "One"})
   performRequest(router, "POST", "/posts", Post{Title: "Second", Content: "Two"})
   // Retrieve list
   w := performRequest(router, "GET", "/posts", nil)
   if w.Code != http.StatusOK {
       t.Fatalf("expected status 200, got %d", w.Code)
   }
   var posts []Post
   if err := json.Unmarshal(w.Body.Bytes(), &posts); err != nil {
       t.Fatalf("invalid JSON response: %v", err)
   }
   if len(posts) != 2 {
       t.Fatalf("expected 2 posts, got %d", len(posts))
   }
   if posts[0].Title != "Second" {
       t.Fatalf("expected most recent post first, got '%s'", posts[0].Title)
   }
}