package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestNewPostController(t *testing.T) {
	c := newPostController(nil)
	if c == nil {
		t.Error("Invalid post controller")
	}
}

func TestViewPostIndex(t *testing.T) {
	c := newPostController(nil)
	if c == nil {
		t.Error("Invalid post controller")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatalf("Failed to construct temp request object: %s", err)
	}

	c.index(w, req)

	if w.Code != http.StatusOK {
		t.Error("Expected a 200, got %d", w.Code)
	}

	for _, event := range postHistory {
		e, ok := event.(postAdded)
		if !ok {
			continue
		}
		matched, err := regexp.Match(e.content.title, w.Body.Bytes())
		if err != nil {
			t.Fatalf("Error generating match for '%s': %s", e.content.title, err)
			return
		}
		if !matched {
			t.Errorf("Expected string '%s' did not match.", e.content.title)
		}
	}
}

func TestViewPostAdd(t *testing.T) {
	c := newPostController(nil)
	if c == nil {
		t.Error("Invalid post controller")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/posts/add", nil)
	if err != nil {
		t.Fatalf("Failed to construct temp request object: %s", err)
	}

	c.viewAdd(w, req)

	if w.Code != http.StatusOK {
		t.Error("Expected a 200, got %d", w.Code)
	}

}
