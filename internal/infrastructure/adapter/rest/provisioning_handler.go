package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"telecomx-provisioning-service/internal/application/service"
	"telecomx-provisioning-service/internal/domain/model"
)

type ProvisioningHandler struct {
	service *service.ProvisioningService
}

func NewProvisioningHandler(s *service.ProvisioningService) *ProvisioningHandler {
	return &ProvisioningHandler{service: s}
}

func (h *ProvisioningHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/provisioning", h.handleProvisioning)
}

func (h *ProvisioningHandler) handleProvisioning(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		data, err := h.service.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)

	case http.MethodPost:
		var p model.Provisioning
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.service.Create(ctx, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
