package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	FlareDefine "github.com/soulteary/flare/config/define"
	FlareMDI "github.com/soulteary/flare/internal/resources/mdi"
)

func TestApiKeyMatch(t *testing.T) {
	if !apiKeyMatch("Bearer secret", "", "secret") {
		t.Fatal("bearer should match")
	}
	if !apiKeyMatch("", "secret", "secret") {
		t.Fatal("x-api-key should match")
	}
	if apiKeyMatch("Bearer wrong", "nope", "secret") {
		t.Fatal("bad keys should fail")
	}
}

func TestSearchIconsWired(t *testing.T) {
	// icons.go must be present from build
	if !FlareMDI.IconExists("home-line") {
		t.Skip("icons not generated")
	}
	names := FlareMDI.SearchIcons("home-line", 5)
	if len(names) == 0 {
		t.Fatal("expected matches")
	}
}

func TestRegisterRoutingDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	FlareDefine.AppFlags.EnableAPI = false
	r := gin.New()
	RegisterRouting(r)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when disabled, got %d", w.Code)
	}
}

func TestRegisterRoutingEnabledIndex(t *testing.T) {
	gin.SetMode(gin.TestMode)
	FlareDefine.AppFlags.EnableAPI = true
	FlareDefine.AppFlags.APIKey = ""
	r := gin.New()
	RegisterRouting(r)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}
