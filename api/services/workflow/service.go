package workflow

import (
	"net/http"
	"workflow-code-test/api/internal/workflow"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	db *pgx.Conn
}

func NewService(db *pgx.Conn) (*Service, error) {
	return &Service{db: db}, nil
}

// jsonMiddleware sets the Content-Type header to application/json
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Service) LoadRoutes(parentRouter *mux.Router, isProduction bool) {
	router := parentRouter.PathPrefix("/workflows").Subrouter()
	router.StrictSlash(false)
	router.Use(jsonMiddleware)

	repo := workflow.NewRepository()
	svc := workflow.NewService(repo)
	wh := workflow.NewHandler(svc)

	router.HandleFunc("/{id}", wh.Workflow).Methods("GET")
	router.HandleFunc("/{id}/execute", s.HandleExecuteWorkflow).Methods("POST")
}
