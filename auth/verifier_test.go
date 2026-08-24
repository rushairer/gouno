package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestVerifierRejectsWrongSigningMethod(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	// HS256 签名的 token：keyfunc 只接受 RSA 方法，应被拒绝（防算法混淆）
	hmacKey := []byte("secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString(hmacKey)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	if _, err := verifier.Verify(signed, Options{}); err == nil {
		t.Fatal("expected token with non-RSA signing method to be rejected")
	}
}

func TestVerifierRejectsExpiredToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	claims := jwt.MapClaims{
		"sub": "user-1",
		"iss": "https://issuer",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	token := signToken(t, key, claims)

	if _, err := verifier.Verify(token, Options{Issuer: "https://issuer"}); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestVerifierRejectsMalformedToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	for _, token := range []string{"not-a-jwt", "", "a.b.c"} {
		if _, err := verifier.Verify(token, Options{}); err == nil {
			t.Errorf("expected malformed token %q to be rejected", token)
		}
	}
}

func TestVerifierRejectsUnknownKid(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	claims := jwt.MapClaims{
		"sub": "user-1",
		"iss": "https://issuer",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "ghost-kid" // 不在 JWKS 中
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	if _, err := verifier.Verify(signed, Options{}); err == nil {
		t.Fatal("expected token with unknown kid to be rejected")
	}
}

func TestVerifierRejectsClientIDMismatch(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	claims := jwt.MapClaims{
		"sub":       "user-1",
		"client_id": "other-app",
		"exp":       time.Now().Add(time.Hour).Unix(),
	}
	token := signToken(t, key, claims)

	if _, err := verifier.Verify(token, Options{ClientID: "blog-spa"}); err == nil {
		t.Fatal("expected client_id mismatch to be rejected")
	}
}

func TestVerifierAcceptsClientIDMatchWithoutAzp(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	verifier := newVerifier(t, &key.PublicKey)

	// 无 azp，仅 client_id 匹配：应通过（azp 为空时跳过校验）
	claims := jwt.MapClaims{
		"sub":       "user-1",
		"client_id": "blog-spa",
		"exp":       time.Now().Add(time.Hour).Unix(),
	}
	token := signToken(t, key, claims)

	if _, err := verifier.Verify(token, Options{ClientID: "blog-spa"}); err != nil {
		t.Fatalf("expected valid client_id to pass: %v", err)
	}
}

func TestVerifierRejectsOversizedJWKS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", 1<<20+1)))
	}))
	defer server.Close()
	verifier := &Verifier{jwksURL: server.URL, keys: make(map[string]*rsa.PublicKey)}
	if err := verifier.refreshKeys(); err == nil {
		t.Fatal("expected oversized JWKS response to fail")
	}
}

func TestVerifierRejectsInvalidJWKSJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("{invalid json"))
	}))
	defer server.Close()
	verifier := &Verifier{jwksURL: server.URL, keys: make(map[string]*rsa.PublicKey)}
	if err := verifier.refreshKeys(); err == nil {
		t.Fatal("expected invalid JWKS JSON to fail")
	}
}

func TestVerifierFiltersNonRS256Keys(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	// JWKS 中只有 alg=RS512 的 key：refreshKeys 应过滤掉，导致签名验证找不到公钥
	jwks := JWKS{Keys: []JWK{{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS512",
		Kid: "test-kid",
		N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(bigEndianIntBytes(key.E)),
	}}}
	body, err := json.Marshal(jwks)
	if err != nil {
		t.Fatalf("marshal jwks: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	defer server.Close()

	verifier := &Verifier{jwksURL: server.URL, keys: make(map[string]*rsa.PublicKey)}
	if err := verifier.refreshKeys(); err != nil {
		t.Fatalf("refreshKeys: %v", err)
	}
	if len(verifier.keys) != 0 {
		t.Errorf("expected non-RS256 keys to be filtered, got %d keys", len(verifier.keys))
	}

	claims := jwt.MapClaims{"sub": "user-1", "exp": time.Now().Add(time.Hour).Unix()}
	token := signToken(t, key, claims)
	if _, err := verifier.Verify(token, Options{}); err == nil {
		t.Fatal("expected verify to fail when no usable key exists")
	}
}

func TestParseJWKInvalidBase64(t *testing.T) {
	// N 或 E 包含非法 base64 字符时应返回错误
	if _, err := parseJWK(JWK{N: "!!!not-base64!!!", E: "AQAB"}); err == nil {
		t.Error("expected error for invalid base64 in N")
	}
	if _, err := parseJWK(JWK{N: "AQAB", E: "!!!not-base64!!!"}); err == nil {
		t.Error("expected error for invalid base64 in E")
	}
}

func TestGetPublicKeyCached(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	// 先手动填充缓存，验证命中缓存时不触发网络请求
	verifier := &Verifier{
		keys: map[string]*rsa.PublicKey{"cached-kid": &key.PublicKey},
	}
	got, err := verifier.GetPublicKey("cached-kid")
	if err != nil {
		t.Fatalf("GetPublicKey: %v", err)
	}
	if got != &key.PublicKey {
		t.Error("expected cached key to be returned")
	}
}

func TestGetPublicKeyNotFound(t *testing.T) {
	// 空 JWKS 且缓存为空：kid 不存在
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"keys":[]}`))
	}))
	defer server.Close()
	verifier := &Verifier{jwksURL: server.URL, keys: make(map[string]*rsa.PublicKey)}
	if _, err := verifier.GetPublicKey("missing-kid"); err == nil {
		t.Fatal("expected error for unknown kid")
	}
}
