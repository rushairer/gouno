package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwkFromPublicKey(key *rsa.PublicKey) JWK {
	return JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: "test-kid",
		N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(bigEndianIntBytes(key.E)),
	}
}

func bigEndianIntBytes(value int) []byte {
	var buf []byte
	for value > 0 {
		buf = append([]byte{byte(value & 0xff)}, buf...)
		value >>= 8
	}
	if len(buf) == 0 {
		buf = []byte{0}
	}
	return buf
}

func signToken(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-kid"
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func newVerifier(t *testing.T, key *rsa.PublicKey) *Verifier {
	t.Helper()
	jwks := JWKS{Keys: []JWK{jwkFromPublicKey(key)}}
	body, err := json.Marshal(jwks)
	if err != nil {
		t.Fatalf("marshal jwks: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return NewVerifier(server.URL)
}

func TestVerifierVerifyValidToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	claims := jwt.MapClaims{
		"sub":   "user-1",
		"iss":   "https://issuer",
		"aud":   "blog-spa",
		"azp":   "blog-spa",
		"roles": []interface{}{"admin"},
		"exp":   time.Now().Add(time.Hour).Unix(),
		"nbf":   time.Now().Add(-time.Minute).Unix(),
		"iat":   time.Now().Add(-time.Minute).Unix(),
	}
	token := signToken(t, key, claims)

	got, err := verifier.Verify(token, Options{Issuer: "https://issuer", Audience: "blog-spa", ClientID: "blog-spa"})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got["sub"] != "user-1" {
		t.Fatalf("unexpected sub: %v", got["sub"])
	}
}

func TestVerifierRejectsWrongIssuer(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	claims := jwt.MapClaims{
		"sub": "user-1",
		"iss": "https://other-issuer",
		"aud": "blog-spa",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := signToken(t, key, claims)

	if _, err := verifier.Verify(token, Options{Issuer: "https://issuer"}); err == nil {
		t.Fatal("expected issuer mismatch error")
	}
}

func TestVerifierRejectsWrongClient(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	claims := jwt.MapClaims{
		"sub": "user-1",
		"aud": "blog-spa",
		"azp": "someone-else",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := signToken(t, key, claims)

	if _, err := verifier.Verify(token, Options{Audience: "blog-spa", ClientID: "blog-spa"}); err == nil {
		t.Fatal("expected authorized party mismatch error")
	}
}

func TestVerifierSingleflightAndCooldown(t *testing.T) {
	requestCount := 0
	var mu sync.Mutex
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		jwks := JWKS{Keys: []JWK{jwkFromPublicKey(&privateKey.PublicKey)}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer server.Close()

	verifier := &Verifier{
		jwksURL: server.URL,
		keys:    make(map[string]*rsa.PublicKey),
	}

	// Trigger concurrent key refreshes.
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = verifier.refreshKeys()
		}()
	}
	wg.Wait()

	mu.Lock()
	countBefore := requestCount
	mu.Unlock()
	if countBefore != 1 {
		t.Errorf("expected exactly 1 HTTP request due to singleflight, got %d", countBefore)
	}

	// A sequential request immediately after should be ignored due to cooldown.
	if err := verifier.refreshKeys(); err != nil {
		t.Fatalf("refreshKeys: %v", err)
	}
	mu.Lock()
	countAfter := requestCount
	mu.Unlock()
	if countAfter != 1 {
		t.Errorf("expected request count to remain 1 due to cooldown, got %d", countAfter)
	}
}

func TestVerifierRejectsNonSuccessJWKSResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	verifier := &Verifier{jwksURL: server.URL, keys: make(map[string]*rsa.PublicKey)}
	if err := verifier.refreshKeys(); err == nil {
		t.Fatal("expected non-success JWKS response to fail")
	}
}
