package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_RouteRecording(t *testing.T) {
	r := New()

	r.HandleFunc("POST /api/auth/login", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.HandleFunc("GET /api/user/profile", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if r.TotalRoutes() != 2 {
		t.Fatalf("Expected 2 routes, got %d", r.TotalRoutes())
	}

	routes := r.Routes()
	if routes[0].Method != "POST" || routes[0].Path != "/api/auth/login" {
		t.Errorf("Route 0 mismatch: %+v", routes[0])
	}
	if routes[1].Method != "GET" || routes[1].Path != "/api/user/profile" {
		t.Errorf("Route 1 mismatch: %+v", routes[1])
	}
}

func TestRequestLogger_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created item"))
	})

	logger := RequestLogger(true)
	wrappedHandler := logger(handler)

	req := httptest.NewRequest("POST", "/api/items", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}
	if rec.Body.String() != "created item" {
		t.Errorf("Expected 'created item', got '%s'", rec.Body.String())
	}
}
