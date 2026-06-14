package gouno_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rushairer/gouno"
)

func TestNewResponse(t *testing.T) {
	resp := gouno.NewResponse(200, "ok", "data")
	if resp.Code != 200 {
		t.Errorf("Code = %d; want 200", resp.Code)
	}
	if resp.Message != "ok" {
		t.Errorf("Message = %s; want ok", resp.Message)
	}
	if resp.Data != "data" {
		t.Errorf("Data = %v; want data", resp.Data)
	}
}

func TestNewSuccessResponse(t *testing.T) {
	t.Run("string data", func(t *testing.T) {
		resp := gouno.NewSuccessResponse("hello")
		if resp.Code != http.StatusOK {
			t.Errorf("Code = %d; want %d", resp.Code, http.StatusOK)
		}
		if resp.Message != "success" {
			t.Errorf("Message = %s; want success", resp.Message)
		}
		if resp.Data != "hello" {
			t.Errorf("Data = %v; want hello", resp.Data)
		}
	})

	t.Run("nil data", func(t *testing.T) {
		resp := gouno.NewSuccessResponse(nil)
		if resp.Code != http.StatusOK {
			t.Errorf("Code = %d; want %d", resp.Code, http.StatusOK)
		}
		if resp.Data != nil {
			t.Errorf("Data should be nil")
		}
	})

	t.Run("struct data", func(t *testing.T) {
		type payload struct {
			Key string `json:"key"`
		}
		resp := gouno.NewSuccessResponse(payload{Key: "value"})
		if resp.Code != http.StatusOK {
			t.Errorf("Code = %d; want %d", resp.Code, http.StatusOK)
		}
		d, ok := resp.Data.(payload)
		if !ok {
			t.Fatalf("Data type = %T; want payload", resp.Data)
		}
		if d.Key != "value" {
			t.Errorf("Data.Key = %s; want value", d.Key)
		}
	})
}

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{"bad request", http.StatusBadRequest, "bad request"},
		{"unauthorized", http.StatusUnauthorized, "unauthorized"},
		{"not found", http.StatusNotFound, "not found"},
		{"internal error", http.StatusInternalServerError, "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := gouno.NewErrorResponse(tt.code, tt.message)
			if resp.Code != tt.code {
				t.Errorf("Code = %d; want %d", resp.Code, tt.code)
			}
			if resp.Message != tt.message {
				t.Errorf("Message = %s; want %s", resp.Message, tt.message)
			}
			if resp.Data != nil {
				t.Errorf("Data = %v; want nil", resp.Data)
			}
		})
	}
}

func TestPresetErrorResponses(t *testing.T) {
	tests := []struct {
		name     string
		resp     *gouno.Response
		wantCode int
	}{
		{"InternalServerError", gouno.InternalServerErrorResponse, http.StatusInternalServerError},
		{"BadRequest", gouno.BadRequestResponse, http.StatusBadRequest},
		{"Unauthorized", gouno.UnauthorizedResponse, http.StatusUnauthorized},
		{"Forbidden", gouno.ForbiddenResponse, http.StatusForbidden},
		{"NotFound", gouno.NotFoundResponse, http.StatusNotFound},
		{"MethodNotAllowed", gouno.MethodNotAllowedResponse, http.StatusMethodNotAllowed},
		{"RequestTimeout", gouno.RequestTimeoutResponse, http.StatusRequestTimeout},
		{"Conflict", gouno.ConflictResponse, http.StatusConflict},
		{"Gone", gouno.GoneResponse, http.StatusGone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Code != tt.wantCode {
				t.Errorf("Code = %d; want %d", tt.resp.Code, tt.wantCode)
			}
			if tt.resp.Data != nil {
				t.Errorf("Data should be nil for error responses")
			}
		})
	}
}

func TestNewErrorResponseFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() *gouno.Response
		wantCode int
	}{
		{"InternalServerError", gouno.NewInternalServerErrorResponse, http.StatusInternalServerError},
		{"BadRequest", gouno.NewBadRequestResponse, http.StatusBadRequest},
		{"Unauthorized", gouno.NewUnauthorizedResponse, http.StatusUnauthorized},
		{"Forbidden", gouno.NewForbiddenResponse, http.StatusForbidden},
		{"NotFound", gouno.NewNotFoundResponse, http.StatusNotFound},
		{"MethodNotAllowed", gouno.NewMethodNotAllowedResponse, http.StatusMethodNotAllowed},
		{"RequestTimeout", gouno.NewRequestTimeoutResponse, http.StatusRequestTimeout},
		{"Conflict", gouno.NewConflictResponse, http.StatusConflict},
		{"Gone", gouno.NewGoneResponse, http.StatusGone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.fn()
			if resp.Code != tt.wantCode {
				t.Errorf("Code = %d; want %d", resp.Code, tt.wantCode)
			}
			if resp.Data != nil {
				t.Errorf("Data should be nil for error responses")
			}
		})
	}

	// Verify each call returns a fresh instance
	t.Run("returns unique instances", func(t *testing.T) {
		a := gouno.NewInternalServerErrorResponse()
		b := gouno.NewInternalServerErrorResponse()
		if a == b {
			t.Error("NewInternalServerErrorResponse should return unique instances")
		}
	})
}

func TestResponseJSON(t *testing.T) {
	t.Run("success with data", func(t *testing.T) {
		resp := gouno.NewSuccessResponse(map[string]string{"key": "value"})
		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if m["code"].(float64) != float64(http.StatusOK) {
			t.Errorf("code = %v; want %d", m["code"], http.StatusOK)
		}
		if m["message"].(string) != "success" {
			t.Errorf("message = %v; want success", m["message"])
		}
	})

	t.Run("success with nil data omits field", func(t *testing.T) {
		resp := gouno.NewSuccessResponse(nil)
		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if _, exists := m["data"]; exists {
			t.Error("data field should be omitted when nil")
		}
	})

	t.Run("error response", func(t *testing.T) {
		resp := gouno.NotFoundResponse
		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if m["code"].(float64) != float64(http.StatusNotFound) {
			t.Errorf("code = %v; want %d", m["code"], http.StatusNotFound)
		}
	})
}
