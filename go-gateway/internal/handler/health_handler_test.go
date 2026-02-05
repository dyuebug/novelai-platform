package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthHandler_Healthz 测试健康检查 API
func TestHealthHandler_Healthz(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "健康检查成功",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"status":  "ok",
				"service": "go-gateway",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHealthHandler()
			router := setupTestRouter()

			router.GET("/healthz", handler.Healthz)

			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody["status"], resp["status"])
			assert.Equal(t, tt.expectedBody["service"], resp["service"])
		})
	}
}
