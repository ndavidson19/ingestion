package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := Auth("test-token")

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"Valid token", "Bearer test-token", http.StatusOK},
		{"Invalid token", "Bearer wrong-token", http.StatusUnauthorized},
		{"Missing token", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()

			authMiddleware(handler).ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
