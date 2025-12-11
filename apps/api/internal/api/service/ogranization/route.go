package ogranization

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.handleCreateOrganization).Methods("POST")
}

func (h *Handler) handleCreateOrganization(w http.ResponseWriter, r *http.Request) {}
