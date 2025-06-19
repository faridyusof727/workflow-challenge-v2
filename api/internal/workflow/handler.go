package workflow

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
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

	_, err := h.svc.Execute(r.Context(), id, &input)
	if err != nil {
		render.Error(w, r, http.StatusBadRequest, err, h.log)
		return
	}

	// Generate current timestamp
	currentTime := time.Now().Format(time.RFC3339)

	executionJSON := fmt.Sprintf(`{
		"executedAt": "%s",
		"status": "failed",
		"steps": [
			{
				"nodeId": "start",
				"type": "start",
				"label": "Start",
				"description": "Begin weather check workflow",
				"status": "completed"
			},
			{
				"nodeId": "form",
				"type": "form",
				"label": "User Input",
				"description": "Process collected data - name, email, location",
				"status": "completed",
				"output": {
					"name": "Alice",
					"email": "alice@example.com",
					"city": "Sydney"
				}
			},
			{
				"nodeId": "weather-api",
				"type": "integration",
				"label": "Weather API",
				"description": "Fetch current temperature for Sydney",
				"status": "completed",
				"output": {
					"temperature": 28.5,
					"location": "Sydney"
				}
			},
			{
				"nodeId": "condition",
				"type": "condition",
				"label": "Check Condition",
				"description": "Evaluate temperature threshold",
				"status": "completed",
				"output": {
					"conditionMet": true,
					"threshold": 25,
					"operator": "greater_than",
					"actualValue": 28.5,
					"message": "Temperature 28.5°C is greater than 25°C - condition met"
				}
			},
			{
				"nodeId": "email",
				"type": "email",
				"label": "Send Alert",
				"description": "Email weather alert notification",
				"status": "completed",
				"output": {
					"emailDraft": {
						"to": "alice@example.com",
						"from": "weather-alerts@example.com",
						"subject": "Weather Alert",
						"body": "Weather alert for Sydney! Temperature is 28.5°C!",
						"timestamp": "2024-01-15T14:30:24.856Z"
					},
					"deliveryStatus": "sent",
					"messageId": "msg_abc123def456",
					"emailSent": true
				}
			},
			{
				"nodeId": "end",
				"type": "end",
				"label": "Complete",
				"description": "Workflow execution finished",
				"status": "completed"
			}
		]
	}`, currentTime)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(executionJSON))
}

// Workflow implements Handler.
func (h *HandlerImpl) Workflow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := uuid.Validate(id); err != nil {
		render.Error(w, r, http.StatusBadRequest, render.ErrInvalidWorkflowID, h.log)
		return
	}

	workflow, err := h.svc.Workflow(r.Context(), id)
	if err != nil {
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
