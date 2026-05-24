package router

import (
	"database/sql"

	"github.com/Application-drop-up/Travellle/internal/handler"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	planuc "github.com/Application-drop-up/Travellle/internal/usecase/plan"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB) *chi.Mux {
	planRepo := persistence.NewPlanRepository(db)
	planHandler := handler.NewPlanHandler(planuc.New(planRepo))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	r.Post("/plans", planHandler.Create)
	r.Get("/plans/{share_token}", planHandler.GetByShareToken)

	return r
}
