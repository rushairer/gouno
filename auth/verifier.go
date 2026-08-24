// Package auth provides OIDC resource-server helpers: retrieval of an
// issuer's JSON Web Key Set and verification of RS256 access tokens.
//
// It is the server-side counterpart of the browser @gosso/client SDK and is
// meant to be embedded in any service that validates Gosso-issued tokens
// without importing the identity provider itself.
package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/singleflight"
)

// JWKS mirrors the JSON Web Key Set document published by the issuer.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK is the subset of a JSON Web Key needed to reconstruct an RSA public key.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Options constrains which tokens a resource server accepts.
type Options struct {
	Issuer   string
	Audience string
	ClientID string
}

// Verifier fetches and caches the issuer's JWKS and verifies RS256 access
// tokens against it. It is safe for concurrent use.
type Verifier struct {
	jwksURL       string
	keys          map[string]*rsa.PublicKey
	mu            sync.RWMutex
	sf            singleflight.Group
	lastRefreshed time.Time
}

// NewVerifier returns a Verifier that loads keys from jwksURL. The initial key
// fetch happens in the background so server startup is not blocked.
func NewVerifier(jwksURL string) *Verifier {
	v := &Verifier{
		jwksURL: jwksURL,
		keys:    make(map[string]*rsa.PublicKey),
	}
	go func() {
		for i := 0; i < 10; i++ {
			if err := v.refreshKeys(); err == nil {
				break
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return v
}

// Verify parses tokenStr as a JWT, enforces issuer/audience/client
// constraints, and returns the decoded claims on success.
func (v *Verifier) Verify(tokenStr string, options Options) (jwt.MapClaims, error) {
	parserOptions := make([]jwt.ParserOption, 0, 2)
	if options.Issuer != "" {
		parserOptions = append(parserOptions, jwt.WithIssuer(options.Issuer))
	}
	if options.Audience != "" {
		parserOptions = append(parserOptions, jwt.WithAudience(options.Audience))
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, _ := token.Header["kid"].(string)
		return v.GetPublicKey(kid)
	}, parserOptions...)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if options.ClientID != "" {
		if azp, _ := claims["azp"].(string); azp != "" && azp != options.ClientID {
			return nil, fmt.Errorf("invalid authorized party")
		}
		if clientID, _ := claims["client_id"].(string); clientID != "" && clientID != options.ClientID {
			return nil, fmt.Errorf("invalid client")
		}
	}

	return claims, nil
}

// GetPublicKey returns the cached public key for kid, refreshing from the
// issuer when the key is unknown.
func (v *Verifier) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	key, ok := v.keys[kid]
	v.mu.RUnlock()
	if ok {
		return key, nil
	}

	if err := v.refreshKeys(); err != nil {
		return nil, err
	}

	v.mu.RLock()
	key, ok = v.keys[kid]
	v.mu.RUnlock()
	if ok {
		return key, nil
	}
	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

func (v *Verifier) refreshKeys() error {
	v.mu.RLock()
	lastRefreshed := v.lastRefreshed
	v.mu.RUnlock()

	// 1 minute cooldown to prevent cache stampede / DoS spamming.
	if time.Since(lastRefreshed) < time.Minute {
		return nil
	}

	_, err, _ := v.sf.Do("refresh", func() (interface{}, error) {
		v.mu.RLock()
		lastRefreshedInner := v.lastRefreshed
		v.mu.RUnlock()
		if time.Since(lastRefreshedInner) < time.Minute {
			return nil, nil
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(v.jwksURL)
		if err != nil {
			return nil, err
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("jwks endpoint returned status %d", resp.StatusCode)
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20+1))
		if err != nil {
			return nil, err
		}
		if len(body) > 1<<20 {
			return nil, fmt.Errorf("jwks response exceeds 1 MiB")
		}

		var jwks JWKS
		if err := json.Unmarshal(body, &jwks); err != nil {
			return nil, err
		}

		newKeys := make(map[string]*rsa.PublicKey)
		for _, key := range jwks.Keys {
			if key.Kty == "RSA" && key.Use == "sig" && key.Alg == "RS256" {
				if pubKey, err := parseJWK(key); err == nil {
					newKeys[key.Kid] = pubKey
				}
			}
		}

		v.mu.Lock()
		v.keys = newKeys
		v.lastRefreshed = time.Now()
		v.mu.Unlock()
		return nil, nil
	})

	return err
}

func parseJWK(key JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	var eVal int
	for _, b := range eBytes {
		eVal = (eVal << 8) | int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eVal,
	}, nil
}
