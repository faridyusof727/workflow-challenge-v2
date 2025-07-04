package workflow

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"workflow-code-test/api/pkg/render"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HandlerImpl struct {
	svc Service
	log *slog.Logger
}

// Execute implements Handler.
func (h *HandlerImpl) Execute(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	h.log.Debug("Handling workflow execution for id", "id", id)

	var input ExecutionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Error(w, r, http.StatusBadRequest, err, h.log)
		return
	}

	executionResult, err := h.svc.Execute(r.Context(), id, &input)
	if err != nil {
		h.log.Error("problem finishing workflow execution", slog.Any("ID", id), slog.Any("ERROR", err))
	}
	if executionResult == nil {
		render.Error(w, r, http.StatusBadRequest, fmt.Errorf("got empty executionResult"), h.log)
		return
	}

	render.JSON(w, r, http.StatusOK, executionResult)
}

// Workflow implements Handler.
func (h *HandlerImpl) Workflow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := uuid.Validate(id); err != nil {
		h.log.Error("problem validating workflow id", slog.Any("ID", id), slog.Any("ERROR", err))
		render.Error(w, r, http.StatusBadRequest, render.ErrInvalidWorkflowID, h.log)
		return
	}

	workflow, err := h.svc.Workflow(r.Context(), id)
	if err != nil {
		h.log.Error("problem fetching workflow", slog.Any("ID", id), slog.Any("ERROR", err))
		render.Error(w, r, http.StatusNotFound, err, h.log)
		return
	}

	render.JSON(w, r, http.StatusOK, workflow)
}

func NewHandler(svc Service, log *slog.Logger) Handler {
	return &HandlerImpl{
		svc: svc,
		log: log,
	}
}
