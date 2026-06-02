package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

func TestRedirectToServerAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldServerAddress := system_setting.ServerAddress
	system_setting.ServerAddress = "https://api.f0.nz"
	t.Cleanup(func() {
		system_setting.ServerAddress = oldServerAddress
	})

	tests := []struct {
		name         string
		host         string
		proto        string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "redirects matching HTTP host",
			host:         "api.f0.nz",
			proto:        "http",
			wantStatus:   http.StatusPermanentRedirect,
			wantLocation: "https://api.f0.nz/api/status?x=1",
		},
		{
			name:       "skips matching HTTPS host",
			host:       "api.f0.nz",
			proto:      "https",
			wantStatus: http.StatusOK,
		},
		{
			name:       "skips local host",
			host:       "127.0.0.1:3005",
			proto:      "http",
			wantStatus: http.StatusOK,
		},
		{
			name:       "skips missing forwarded proto",
			host:       "api.f0.nz",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RedirectToServerAddress())
			router.GET("/api/status", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/status?x=1", nil)
			req.Host = tt.host
			if tt.proto != "" {
				req.Header.Set("X-Forwarded-Proto", tt.proto)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if location := recorder.Header().Get("Location"); location != tt.wantLocation {
				t.Fatalf("location = %q, want %q", location, tt.wantLocation)
			}
		})
	}
}
