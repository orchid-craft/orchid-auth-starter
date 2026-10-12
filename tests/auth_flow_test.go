package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yudhiana/web-api/internal/bootstrap/container"
	restfulServer "github.com/yudhiana/web-api/internal/bootstrap/server/restful_server"
	"github.com/yudhiana/web-api/pkg/cookie"
)

func setupTestApp() http.Handler {
	c, err := container.NewContainer()
	if err != nil {
		panic(err)
	}

	restfulServer.NewServer(c)
	return c.GetApp()
}

func TestEndToEndAuthFlow(t *testing.T) {
	app := setupTestApp()

	// 1. Registrasi Pengguna
	regPayload := `{"email":"budi@example.com","password":"password123","display_name":"Budi Santoso","role":"user"}`
	regReq := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(regPayload))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	app.ServeHTTP(regRec, regReq)

	if regRec.Code != http.StatusOK {
		t.Fatalf("Registrasi gagal, status: %d, body: %s", regRec.Code, regRec.Body.String())
	}

	// 2. Login Pengguna
	loginPayload := `{"email":"budi@example.com","password":"password123"}`
	loginReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	app.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("Login gagal, status: %d, body: %s", loginRec.Code, loginRec.Body.String())
	}

	// Verifikasi Cookie Atribut (HttpOnly, SameSite, MaxAge)
	cookies := loginRec.Result().Cookies()
	var accessTokenCookie, refreshTokenCookie, csrfTokenCookie *http.Cookie
	for _, c := range cookies {
		switch c.Name {
		case cookie.AccessTokenCookieName:
			accessTokenCookie = c
		case cookie.RefreshTokenCookieName:
			refreshTokenCookie = c
		case cookie.CSRFCookieName:
			csrfTokenCookie = c
		}
	}

	if accessTokenCookie == nil {
		t.Fatal("Cookie access_token tidak ditemukan")
	}
	if !accessTokenCookie.HttpOnly {
		t.Errorf("access_token HARUS memiliki HttpOnly=true untuk mencegah XSS")
	}
	if accessTokenCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("access_token diharapkan SameSite=Lax, didapat %v", accessTokenCookie.SameSite)
	}

	if refreshTokenCookie == nil {
		t.Fatal("Cookie refresh_token tidak ditemukan")
	}
	if !refreshTokenCookie.HttpOnly {
		t.Errorf("refresh_token HARUS memiliki HttpOnly=true untuk mencegah XSS")
	}
	if refreshTokenCookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("refresh_token diharapkan SameSite=Strict, didapat %v", refreshTokenCookie.SameSite)
	}
	if refreshTokenCookie.Path != "/api/auth" {
		t.Errorf("refresh_token diharapkan memiliki path restricted /api/auth, didapat %s", refreshTokenCookie.Path)
	}

	if csrfTokenCookie == nil {
		t.Fatal("Cookie csrf_token tidak ditemukan")
	}
	if csrfTokenCookie.HttpOnly {
		t.Errorf("csrf_token HARUS HttpOnly=false agar dapat dibaca oleh script Frontend (Double Submit)")
	}

	// 3. Akses Protected Endpoint (GET /api/me) dengan Access Token Cookie
	profileReq := httptest.NewRequest("GET", "/api/me", nil)
	profileReq.AddCookie(accessTokenCookie)
	profileRec := httptest.NewRecorder()
	app.ServeHTTP(profileRec, profileReq)

	if profileRec.Code != http.StatusOK {
		t.Fatalf("Gagal mengakses protected profile, status: %d, body: %s", profileRec.Code, profileRec.Body.String())
	}

	// 4. Test Proteksi CSRF pada Mutasi Data (PUT /api/me)
	// Skenario 4a: Tanpa header X-CSRF-Token -> HARUS DITOLAK (Status 403 Forbidden)
	putReqWithoutCSRF := httptest.NewRequest("PUT", "/api/me", bytes.NewBufferString(`{"display_name":"Budi Update"}`))
	putReqWithoutCSRF.Header.Set("Content-Type", "application/json")
	putReqWithoutCSRF.AddCookie(accessTokenCookie)
	putReqWithoutCSRF.AddCookie(csrfTokenCookie)
	putRecWithoutCSRF := httptest.NewRecorder()
	app.ServeHTTP(putRecWithoutCSRF, putReqWithoutCSRF)

	if putRecWithoutCSRF.Code != http.StatusForbidden {
		t.Fatalf("Diharapkan 403 Forbidden ketika X-CSRF-Token absen, tapi didapat: %d, body: %s", putRecWithoutCSRF.Code, putRecWithoutCSRF.Body.String())
	}

	// Skenario 4b: Dengan X-CSRF-Token yang salah/palsu -> HARUS DITOLAK (Status 403 Forbidden)
	putReqFakeCSRF := httptest.NewRequest("PUT", "/api/me", bytes.NewBufferString(`{"display_name":"Budi Update"}`))
	putReqFakeCSRF.Header.Set("Content-Type", "application/json")
	putReqFakeCSRF.Header.Set("X-CSRF-Token", "fake-tampered-token")
	putReqFakeCSRF.AddCookie(accessTokenCookie)
	putReqFakeCSRF.AddCookie(csrfTokenCookie)
	putRecFakeCSRF := httptest.NewRecorder()
	app.ServeHTTP(putRecFakeCSRF, putReqFakeCSRF)

	if putRecFakeCSRF.Code != http.StatusForbidden {
		t.Fatalf("Diharapkan 403 Forbidden ketika X-CSRF-Token palsu, tapi didapat: %d", putRecFakeCSRF.Code)
	}

	// Skenario 4c: Dengan X-CSRF-Token yang valid sesuai cookie -> HARUS SUKSES (Status 200 OK)
	putReqValid := httptest.NewRequest("PUT", "/api/me", bytes.NewBufferString(`{"display_name":"Budi Santoso Baru"}`))
	putReqValid.Header.Set("Content-Type", "application/json")
	putReqValid.Header.Set("X-CSRF-Token", csrfTokenCookie.Value)
	putReqValid.AddCookie(accessTokenCookie)
	putReqValid.AddCookie(csrfTokenCookie)
	putRecValid := httptest.NewRecorder()
	app.ServeHTTP(putRecValid, putReqValid)

	if putRecValid.Code != http.StatusOK {
		t.Fatalf("Diharapkan 200 OK ketika X-CSRF-Token valid, tapi didapat: %d, body: %s", putRecValid.Code, putRecValid.Body.String())
	}

	// 5. Test Refresh Token Rotation (RTR) (POST /api/auth/refresh)
	refreshReq := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	refreshReq.Header.Set("X-CSRF-Token", csrfTokenCookie.Value)
	refreshReq.AddCookie(refreshTokenCookie)
	refreshReq.AddCookie(csrfTokenCookie)
	refreshRec := httptest.NewRecorder()
	app.ServeHTTP(refreshRec, refreshReq)

	if refreshRec.Code != http.StatusOK {
		t.Fatalf("Gagal melakukan refresh token, status: %d, body: %s", refreshRec.Code, refreshRec.Body.String())
	}

	// Dapatkan cookie baru hasil rotasi
	newCookies := refreshRec.Result().Cookies()
	var newRefreshTokenCookie, newAccessTokenCookie *http.Cookie
	for _, c := range newCookies {
		if c.Name == cookie.RefreshTokenCookieName {
			newRefreshTokenCookie = c
		}
		if c.Name == cookie.AccessTokenCookieName {
			newAccessTokenCookie = c
		}
	}
	if newRefreshTokenCookie == nil || newAccessTokenCookie == nil {
		t.Fatal("Pasangan token baru hasil rotasi tidak ditemukan di cookie")
	}

	if newRefreshTokenCookie.Value == refreshTokenCookie.Value {
		t.Errorf("Refresh token lama tidak dirotasi! Nilai token sama persis")
	}

	// 6. DETEKSI REUSE ATTACK: Coba gunakan Refresh Token LAMA yang sudah di-rotasi
	reuseReq := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	reuseReq.Header.Set("X-CSRF-Token", csrfTokenCookie.Value)
	reuseReq.AddCookie(refreshTokenCookie) // Gunakan token lama!
	reuseReq.AddCookie(csrfTokenCookie)
	reuseRec := httptest.NewRecorder()
	app.ServeHTTP(reuseRec, reuseReq)

	if reuseRec.Code != http.StatusForbidden {
		t.Fatalf("Diharapkan 403 Forbidden saat reuse attack terjadi, tetapi didapat: %d, body: %s", reuseRec.Code, reuseRec.Body.String())
	}

	// 7. Logout dan Pembersihan Cookie
	logoutReq := httptest.NewRequest("POST", "/api/auth/logout", nil)
	logoutReq.Header.Set("X-CSRF-Token", csrfTokenCookie.Value)
	logoutReq.AddCookie(newAccessTokenCookie)
	logoutReq.AddCookie(newRefreshTokenCookie)
	logoutReq.AddCookie(csrfTokenCookie)
	logoutRec := httptest.NewRecorder()
	app.ServeHTTP(logoutRec, logoutReq)

	if logoutRec.Code != http.StatusOK {
		t.Fatalf("Logout gagal, status: %d, body: %s", logoutRec.Code, logoutRec.Body.String())
	}

	// Periksa apakah cookies dibersihkan (MaxAge: -1)
	for _, c := range logoutRec.Result().Cookies() {
		if c.MaxAge != -1 {
			t.Errorf("Cookie %s harus memiliki MaxAge=-1 saat logout, didapat: %d", c.Name, c.MaxAge)
		}
	}
}

func TestRoleBasedAccessControl(t *testing.T) {
	app := setupTestApp()

	// 1. Login sebagai user biasa (Bob Member)
	bobLoginPayload := `{"email":"bob@example.com","password":"password123"}`
	bobReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(bobLoginPayload))
	bobReq.Header.Set("Content-Type", "application/json")
	bobRec := httptest.NewRecorder()
	app.ServeHTTP(bobRec, bobReq)

	if bobRec.Code != http.StatusOK {
		t.Fatalf("Login Bob gagal: %d, body: %s", bobRec.Code, bobRec.Body.String())
	}

	var bobAccessToken *http.Cookie
	for _, c := range bobRec.Result().Cookies() {
		if c.Name == cookie.AccessTokenCookieName {
			bobAccessToken = c
		}
	}
	if bobAccessToken == nil {
		t.Fatal("Access token Bob tidak ditemukan")
	}

	// 2. Akses GET /api/users sebagai user biasa -> HARUS 403 Forbidden (P0-2)
	listUsersReq := httptest.NewRequest("GET", "/api/users", nil)
	listUsersReq.AddCookie(bobAccessToken)
	listUsersRec := httptest.NewRecorder()
	app.ServeHTTP(listUsersRec, listUsersReq)

	if listUsersRec.Code != http.StatusForbidden {
		t.Fatalf("Diharapkan 403 Forbidden untuk user biasa pada /api/users, didapat: %d, body: %s", listUsersRec.Code, listUsersRec.Body.String())
	}

	// 3. Login sebagai Admin (Alice Admin)
	aliceLoginPayload := `{"email":"alice@example.com","password":"password123"}`
	aliceReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(aliceLoginPayload))
	aliceReq.Header.Set("Content-Type", "application/json")
	aliceRec := httptest.NewRecorder()
	app.ServeHTTP(aliceRec, aliceReq)

	if aliceRec.Code != http.StatusOK {
		t.Fatalf("Login Alice gagal: %d, body: %s", aliceRec.Code, aliceRec.Body.String())
	}

	var aliceAccessToken *http.Cookie
	for _, c := range aliceRec.Result().Cookies() {
		if c.Name == cookie.AccessTokenCookieName {
			aliceAccessToken = c
		}
	}
	if aliceAccessToken == nil {
		t.Fatal("Access token Alice tidak ditemukan")
	}

	// 4. Akses GET /api/users sebagai Admin -> HARUS 200 OK
	adminListUsersReq := httptest.NewRequest("GET", "/api/users", nil)
	adminListUsersReq.AddCookie(aliceAccessToken)
	adminListUsersRec := httptest.NewRecorder()
	app.ServeHTTP(adminListUsersRec, adminListUsersReq)

	if adminListUsersRec.Code != http.StatusOK {
		t.Fatalf("Diharapkan 200 OK untuk admin pada /api/users, didapat: %d, body: %s", adminListUsersRec.Code, adminListUsersRec.Body.String())
	}
}

func TestInputValidationAndConflict(t *testing.T) {
	app := setupTestApp()

	// 1. Password terlalu pendek (< 8 karakter) -> 400 Bad Request (P0-3)
	shortPassPayload := `{"email":"test1@example.com","password":"short","display_name":"Test User"}`
	req1 := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(shortPassPayload))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	app.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusBadRequest {
		t.Fatalf("Diharapkan 400 Bad Request untuk password pendek, didapat: %d, body: %s", rec1.Code, rec1.Body.String())
	}

	// 2. Format email tidak valid -> 400 Bad Request (P0-3)
	badEmailPayload := `{"email":"not-an-email","password":"password123","display_name":"Test User"}`
	req2 := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(badEmailPayload))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	app.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("Diharapkan 400 Bad Request untuk email tidak valid, didapat: %d, body: %s", rec2.Code, rec2.Body.String())
	}

	// 3. Registrasi valid
	validPayload := `{"email":"charlie@example.com","password":"password123","display_name":"Charlie","role":"admin"}`
	req3 := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(validPayload))
	req3.Header.Set("Content-Type", "application/json")
	rec3 := httptest.NewRecorder()
	app.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Fatalf("Registrasi valid gagal: %d, body: %s", rec3.Code, rec3.Body.String())
	}

	// Verifikasi P0-1: Meskipun payload mencoba role "admin", hasil registrasi tetap "user"
	if !bytes.Contains(rec3.Body.Bytes(), []byte(`"role":"user"`)) {
		t.Fatalf("Pencegahan Privilege Escalation gagal: role pengguna publik harus 'user', body: %s", rec3.Body.String())
	}

	// 4. Registrasi ulang email yang sama -> HARUS 409 Conflict (P2-3)
	req4 := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(validPayload))
	req4.Header.Set("Content-Type", "application/json")
	rec4 := httptest.NewRecorder()
	app.ServeHTTP(rec4, req4)

	if rec4.Code != http.StatusConflict {
		t.Fatalf("Diharapkan 409 Conflict untuk email yang sudah terdaftar, didapat: %d, body: %s", rec4.Code, rec4.Body.String())
	}
}

func TestStrictCookieRefresh(t *testing.T) {
	app := setupTestApp()

	// 1. Dapatkan CSRF token awal
	csrfReq := httptest.NewRequest("GET", "/api/auth/csrf-token", nil)
	csrfRec := httptest.NewRecorder()
	app.ServeHTTP(csrfRec, csrfReq)

	var csrfCookie *http.Cookie
	for _, c := range csrfRec.Result().Cookies() {
		if c.Name == cookie.CSRFCookieName {
			csrfCookie = c
		}
	}
	if csrfCookie == nil {
		t.Fatal("CSRF cookie tidak ditemukan")
	}

	// 2. Coba refresh token hanya lewat JSON body tanpa Cookie -> HARUS 401 Unauthorized (P1-1)
	fallbackBody := `{"refresh_token":"some-token-string"}`
	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBufferString(fallbackBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfCookie.Value)
	req.AddCookie(csrfCookie)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Diharapkan 401 Unauthorized saat refresh token tidak ada di cookie, didapat: %d, body: %s", rec.Code, rec.Body.String())
	}
}
