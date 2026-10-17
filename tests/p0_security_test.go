package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/orchid-craft/orchid-auth-starter/pkg/cookie"
)

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// assertNoTokensInBody memastikan access/refresh token hanya ada di cookie, tidak di body JSON.
func assertNoTokensInBody(t *testing.T, label string, rec *httptest.ResponseRecorder) {
	t.Helper()

	body := rec.Body.String()
	for _, key := range []string{`"access_token"`, `"refresh_token"`} {
		if strings.Contains(body, key) {
			t.Errorf("[%s] body JSON tidak boleh memuat field %s, body: %s", label, key, body)
		}
	}

	for _, c := range rec.Result().Cookies() {
		isAuthToken := c.Name == cookie.AccessTokenCookieName || c.Name == cookie.RefreshTokenCookieName
		if isAuthToken && c.Value != "" && strings.Contains(body, c.Value) {
			t.Errorf("[%s] nilai cookie %s bocor ke body JSON", label, c.Name)
		}
	}
}

func TestTokensNotExposedInResponseBody(t *testing.T) {
	app := setupTestApp()

	// Register
	regReq := httptest.NewRequest("POST", "/api/auth/register",
		bytes.NewBufferString(`{"email":"dewi@example.com","password":"password123","display_name":"Dewi Lestari"}`))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	app.ServeHTTP(regRec, regReq)
	if regRec.Code != http.StatusOK {
		t.Fatalf("register gagal: %d, body: %s", regRec.Code, regRec.Body.String())
	}
	assertNoTokensInBody(t, "register", regRec)
	if c := findCookie(regRec.Result().Cookies(), cookie.AccessTokenCookieName); c == nil || c.Value == "" {
		t.Fatal("register seharusnya tetap menyetel cookie access_token")
	}

	// Login
	loginReq := httptest.NewRequest("POST", "/api/auth/login",
		bytes.NewBufferString(`{"email":"bob@example.com","password":"password123"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	app.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login gagal: %d, body: %s", loginRec.Code, loginRec.Body.String())
	}
	assertNoTokensInBody(t, "login", loginRec)

	refreshCookie := findCookie(loginRec.Result().Cookies(), cookie.RefreshTokenCookieName)
	csrfCookie := findCookie(loginRec.Result().Cookies(), cookie.CSRFCookieName)
	if refreshCookie == nil || refreshCookie.Value == "" || csrfCookie == nil {
		t.Fatal("login seharusnya menyetel cookie refresh_token dan csrf_token")
	}

	// Refresh
	refreshReq := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	refreshReq.Header.Set("X-CSRF-Token", csrfCookie.Value)
	refreshReq.AddCookie(refreshCookie)
	refreshReq.AddCookie(csrfCookie)
	refreshRec := httptest.NewRecorder()
	app.ServeHTTP(refreshRec, refreshReq)
	if refreshRec.Code != http.StatusOK {
		t.Fatalf("refresh gagal: %d, body: %s", refreshRec.Code, refreshRec.Body.String())
	}
	assertNoTokensInBody(t, "refresh", refreshRec)
}

// Dua (atau lebih) request paralel dengan refresh token yang sama tidak boleh sama-sama sukses.
// Jalankan dengan: go test -race ./tests/...
func TestConcurrentRefreshOnlyOneSucceeds(t *testing.T) {
	app := setupTestApp()

	loginReq := httptest.NewRequest("POST", "/api/auth/login",
		bytes.NewBufferString(`{"email":"bob@example.com","password":"password123"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	app.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login gagal: %d, body: %s", loginRec.Code, loginRec.Body.String())
	}

	refreshCookie := findCookie(loginRec.Result().Cookies(), cookie.RefreshTokenCookieName)
	csrfCookie := findCookie(loginRec.Result().Cookies(), cookie.CSRFCookieName)
	if refreshCookie == nil || csrfCookie == nil {
		t.Fatal("cookie refresh_token / csrf_token tidak ditemukan")
	}

	const workers = 20
	statuses := make([]int, workers)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
			req.Header.Set("X-CSRF-Token", csrfCookie.Value)
			req.AddCookie(refreshCookie)
			req.AddCookie(csrfCookie)
			// IP berbeda per request agar rate limiter tidak ikut campur dalam pengujian ini.
			req.RemoteAddr = fmt.Sprintf("198.51.100.%d:40000", i+1)

			rec := httptest.NewRecorder()
			<-start
			app.ServeHTTP(rec, req)
			statuses[i] = rec.Code
		}(i)
	}

	close(start)
	wg.Wait()

	success := 0
	for i, code := range statuses {
		switch code {
		case http.StatusOK:
			success++
		case http.StatusForbidden:
			// reuse terdeteksi: perilaku yang benar bagi request yang kalah
		default:
			t.Errorf("request #%d: status tak terduga %d (seharusnya 200 atau 403)", i, code)
		}
	}

	if success != 1 {
		t.Fatalf("diharapkan tepat 1 refresh sukses dari %d request paralel, didapat %d (statuses: %v)", workers, success, statuses)
	}
}
