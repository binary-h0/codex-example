package main

import (
   "encoding/json"
   "net/http"
   "testing"
)

// TestGetPostInvalidID ensures requesting a post with non-numeric ID returns 404
func TestGetPostInvalidID(t *testing.T) {
   db := setupTestDB(t)
   router := setupRouter(db)
   w := performRequest(router, "GET", "/posts/abc", nil)
   if w.Code != http.StatusNotFound {
       t.Fatalf("expected status 404 for invalid ID, got %d", w.Code)
   }
   var resp map[string]string
   if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
       t.Fatalf("invalid JSON response: %v", err)
   }
   if resp["error"] != "post not found" {
       t.Fatalf("expected error 'post not found', got '%s'", resp["error"])
   }
}