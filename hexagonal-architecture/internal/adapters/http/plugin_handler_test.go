package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/application"
)

func TestPluginHandlerRegistersEnablesAndListsPlugins(t *testing.T) {
	pluginRepo := memory.NewPluginRepository()
	handler := NewPluginHandler(
		application.NewRegisterPricingPluginUseCase(pluginRepo),
		application.NewEnablePluginUseCase(pluginRepo),
		application.NewListPluginsUseCase(pluginRepo),
	)

	registerRequest := httptest.NewRequest(http.MethodPost, "/plugins", strings.NewReader(`{"key":"seasonal-pricing","discountPercent":5}`))
	registerRecorder := httptest.NewRecorder()
	handler.ServeHTTP(registerRecorder, registerRequest)

	if registerRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, registerRecorder.Code)
	}

	enableRequest := httptest.NewRequest(http.MethodPost, "/plugins/seasonal-pricing/enable", nil)
	enableRecorder := httptest.NewRecorder()
	handler.ServeHTTP(enableRecorder, enableRequest)

	if enableRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, enableRecorder.Code)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/plugins", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRecorder.Code)
	}

	body := listRecorder.Body.String()
	if !strings.Contains(body, `"key":"seasonal-pricing"`) || !strings.Contains(body, `"status":"Enabled"`) {
		t.Fatalf("expected enabled plugin in list, got %s", body)
	}
}
