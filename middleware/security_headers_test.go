package middleware

import (
	"net/http"
	"net/http/httptest"
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
