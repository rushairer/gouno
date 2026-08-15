package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateCSRFToken returns a hex-encoded 256-bit random token.
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate csrf token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// CSRFMatches reports whether the submitted token matches the cookie value,
// using a constant-time comparison to avoid timing side channels.
func CSRFMatches(cookie, submitted string) bool {
	if cookie == "" || submitted == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie), []byte(submitted)) == 1
}

// SetCSRFCookie writes a double-submit CSRF cookie. The cookie is deliberately
// HttpOnly=false because the SPA must read it for the double-submit pattern,
// and SameSite=Lax so top-level OAuth redirects from the identity provider are
// not blocked.
func SetCSRFCookie(ctx *gin.Context, name, value string, secure bool, maxAge time.Duration) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: false,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// EnsureCSRFCookie sets the double-submit CSRF cookie, generating a fresh
// token when none is present. It returns an error when token generation fails;
// callers should then abort with their preferred error response.
func EnsureCSRFCookie(ctx *gin.Context, name string, secure bool, maxAge time.Duration) error {
	if cookie, err := ctx.Cookie(name); err == nil && cookie != "" {
		return nil
	}
	token, err := GenerateCSRFToken()
	if err != nil {
		return err
	}
	SetCSRFCookie(ctx, name, token, secure, maxAge)
	return nil
}
