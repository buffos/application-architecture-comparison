package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type PluginHandler struct {
	register application.RegisterPricingPluginUseCase
	enable   application.EnablePluginUseCase
	list     application.ListPluginsUseCase
}

type pluginResponse struct {
	Key             string `json:"key"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	DiscountPercent int    `json:"discountPercent"`
}

func NewPluginHandler(register application.RegisterPricingPluginUseCase, enable application.EnablePluginUseCase, list application.ListPluginsUseCase) PluginHandler {
	return PluginHandler{
		register: register,
		enable:   enable,
		list:     list,
	}
}

func (h PluginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/plugins":
		h.registerPlugin(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/plugins":
		h.listPlugins(w, r)
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/enable"):
		h.enablePlugin(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h PluginHandler) registerPlugin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key             string `json:"key"`
		DiscountPercent int    `json:"discountPercent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plugin, err := h.register.Execute(body.Key, body.DiscountPercent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writePluginResponse(w, http.StatusCreated, plugin)
}

func (h PluginHandler) listPlugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := h.list.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]pluginResponse, 0, len(plugins))
	for _, plugin := range plugins {
		response = append(response, toPluginResponse(plugin))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h PluginHandler) enablePlugin(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/plugins/"), "/enable")

	plugin, err := h.enable.Execute(key)
	if err != nil {
		status := http.StatusBadRequest
		if err == domain.ErrPluginNotFound {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	writePluginResponse(w, http.StatusOK, plugin)
}

func writePluginResponse(w http.ResponseWriter, status int, plugin domain.PluginRegistration) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(toPluginResponse(plugin))
}

func toPluginResponse(plugin domain.PluginRegistration) pluginResponse {
	return pluginResponse{
		Key:             plugin.Key,
		Type:            plugin.Type,
		Status:          plugin.Status,
		DiscountPercent: plugin.DiscountPercent,
	}
}
