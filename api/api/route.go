package api

import (
	"net/http"
	"workflow-code-test/api/internal/workflow"
	"workflow-code-test/api/pkg/di"

	"github.com/gorilla/mux"
)

type Service struct {
	di *di.Container
}

func NewRouter(di *di.Container) (*Service, error) {
	return &Service{
		di: di,
	}, nil
}

func (s *Service) LoadRoutes(parentRouter *mux.Router, isProduction bool) {
	router := parentRouter.PathPrefix("/workflows").Subrouter()
	router.StrictSlash(false)
	router.Use(JsonMiddleware)

	repo := workflow.NewRepository(s.di.DbService.Conn())
	svc := workflow.NewService(repo)
	wh := workflow.NewHandler(svc, s.di.Logger)

	router.HandleFunc("/{id}", wh.Workflow).Methods(http.MethodGet)
	router.HandleFunc("/{id}/execute", wh.Execute).Methods(http.MethodPost)
}
