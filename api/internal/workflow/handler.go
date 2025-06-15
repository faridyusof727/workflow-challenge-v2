package workflow

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
)

type HandlerImpl struct {
	svc Service
}

// Execute implements Handler.
func (h *HandlerImpl) Execute(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// Workflow implements Handler.
func (h *HandlerImpl) Workflow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	workflow, err := h.svc.Workflow(r.Context(), id)
	// TODO: Proper error handling
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, workflow)
}

func NewHandler(svc Service) Handler {
	return &HandlerImpl{
		svc: svc,
	}
}
