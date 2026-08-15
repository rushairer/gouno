package middleware

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersOptions configures the security headers applied to every
// response. CSP and Permissions-Policy are left to the caller because they are
// service-specific; the shared static headers are always set.
type SecurityHeadersOptions struct {
	IsProduction              bool
	CSP                       string // full Content-Security-Policy value; empty disables it
	PermissionsPolicy         string // full Permissions-Policy value; empty disables it
	CrossOriginOpenerPolicy   string // e.g. "same-origin"; empty disables it
	CrossOriginResourcePolicy string // e.g. "same-origin"; empty disables it
}

// SecurityHeaders sets the common browser security headers shared by Gosso
// services. HSTS is only emitted when IsProduction is true (meaningless over
// plain HTTP).
func SecurityHeaders(opts SecurityHeadersOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "0")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		if opts.IsProduction {
			ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		if opts.PermissionsPolicy != "" {
			ctx.Header("Permissions-Policy", opts.PermissionsPolicy)
		}
		if opts.CSP != "" {
			ctx.Header("Content-Security-Policy", opts.CSP)
		}
		if opts.CrossOriginOpenerPolicy != "" {
			ctx.Header("Cross-Origin-Opener-Policy", opts.CrossOriginOpenerPolicy)
		}
		if opts.CrossOriginResourcePolicy != "" {
			ctx.Header("Cross-Origin-Resource-Policy", opts.CrossOriginResourcePolicy)
		}
		ctx.Next()
	}
}

// GenerateCSPNonce returns a base64url-encoded 128-bit random nonce for use in
// a Content-Security-Policy script/style nonce.
func GenerateCSPNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
