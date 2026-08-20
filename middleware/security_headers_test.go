package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeaders(SecurityHeadersOptions{
		IsProduction:              true,
		CSP:                       "default-src 'self'",
		PermissionsPolicy:         "camera=()",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
	}))
	router.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	headers := rec.Header()
	if headers.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing nosniff: %v", headers)
	}
	if headers.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("missing frame deny: %v", headers)
	}
	if headers.Get("Strict-Transport-Security") == "" {
		t.Fatal("expected HSTS in production")
	}
	if headers.Get("Content-Security-Policy") != "default-src 'self'" {
		t.Fatalf("unexpected CSP: %q", headers.Get("Content-Security-Policy"))
	}
	if headers.Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Fatalf("unexpected COOP: %q", headers.Get("Cross-Origin-Opener-Policy"))
	}
	if headers.Get("Cross-Origin-Resource-Policy") != "same-origin" {
		t.Fatalf("unexpected CORP: %q", headers.Get("Cross-Origin-Resource-Policy"))
	}
	if headers.Get("Permissions-Policy") != "camera=()" {
		t.Fatalf("unexpected Permissions-Policy: %q", headers.Get("Permissions-Policy"))
	}
	if headers.Get("X-XSS-Protection") != "0" {
		t.Fatalf("unexpected X-XSS-Protection: %q", headers.Get("X-XSS-Protection"))
	}
	if headers.Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
		t.Fatalf("unexpected Referrer-Policy: %q", headers.Get("Referrer-Policy"))
	}
}

func TestGenerateCSPNonce(t *testing.T) {
	a, err := GenerateCSPNonce()
	if err != nil {
		t.Fatalf("generate nonce: %v", err)
	}
	b, err := GenerateCSPNonce()
	if err != nil {
		t.Fatalf("generate nonce: %v", err)
	}
	// 16 字节随机数经 base64url 编码：固定 22 字符、无填充
	if a == "" || len(a) != 22 {
		t.Fatalf("unexpected nonce length: %d", len(a))
	}
	if a == b {
		t.Fatal("expected distinct nonces")
	}
	if strings.ContainsAny(a, "+/=") {
		t.Fatalf("nonce contains non-urlsafe chars: %q", a)
	}
}

func TestSecurityHeadersNoHSTSWhenNotProduction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeaders(SecurityHeadersOptions{}))
	router.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("did not expect HSTS outside production")
	}
	if rec.Header().Get("Content-Security-Policy") != "" {
		t.Fatal("did not expect CSP when unset")
	}
}
