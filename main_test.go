package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_NotFound(t *testing.T) {
	mux := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d",
			http.StatusNotFound, recorder.Code)
	}

}

func TestRouter_Home(t *testing.T) {
	mux := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d",
			http.StatusOK, recorder.Code)
	}
}

func TestRouter_GenerateRejectsGET(t *testing.T) {
	mux := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d",
			http.StatusMethodNotAllowed, recorder.Code)
	}
}
