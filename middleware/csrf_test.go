package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGenerateCSRFToken(t *testing.T) {
	a, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	b, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if a == "" || len(a) != 64 {
		t.Fatalf("unexpected token length: %d", len(a))
	}
	if a == b {
		t.Fatal("expected distinct tokens")
	}
}

func TestCSRFMatches(t *testing.T) {
	if !CSRFMatches("abc", "abc") {
		t.Fatal("expected match")
	}
	if CSRFMatches("abc", "def") {
		t.Fatal("expected mismatch")
	}
	if CSRFMatches("", "abc") {
		t.Fatal("empty cookie must not match")
	}
	if CSRFMatches("abc", "") {
		t.Fatal("empty header must not match")
	}
}

func TestEnsureCSRFCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", func(ctx *gin.Context) {
		if err := EnsureCSRFCookie(ctx, "csrf", true, time.Hour); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(rec, req)

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != "csrf" || cookie.Value == "" {
		t.Fatalf("unexpected cookie: %+v", cookie)
	}
	if cookie.Secure != true || cookie.HttpOnly != false || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected cookie attributes: %+v", cookie)
	}
}

func TestEnsureCSRFCookieKeepsExisting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	var cookieCount int
	router.GET("/", func(ctx *gin.Context) {
		if err := EnsureCSRFCookie(ctx, "csrf", true, time.Hour); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		cookieCount = len(ctx.Writer.Header().Values("Set-Cookie"))
		ctx.Status(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "csrf", Value: "existing-token"})
	router.ServeHTTP(rec, req)

	if cookieCount != 0 {
		t.Fatalf("expected no new cookie when one already exists, got %d", cookieCount)
	}
}

func TestSetCSRFCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	SetCSRFCookie(ctx, "csrf", "token-123", false, 30*time.Minute)

	cookie := rec.Result().Cookies()
	if len(cookie) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookie))
	}
	c := cookie[0]
	if c.Name != "csrf" || c.Value != "token-123" {
		t.Fatalf("unexpected cookie: %+v", c)
	}
	if c.MaxAge != 1800 { // 30 分钟 → 1800 秒
		t.Fatalf("expected MaxAge 1800, got %d", c.MaxAge)
	}
	if c.Path != "/" || c.Secure || c.HttpOnly {
		t.Fatalf("unexpected attributes: %+v", c)
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", c.SameSite)
	}
}
