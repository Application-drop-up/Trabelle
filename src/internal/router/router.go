package router

import (
	"database/sql"

	"github.com/Application-drop-up/Travellle/internal/handler"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	pinuc "github.com/Application-drop-up/Travellle/internal/usecase/pin"
	planuc "github.com/Application-drop-up/Travellle/internal/usecase/plan"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB) *chi.Mux {
	planHandler := handler.NewPlanHandler(planuc.New(persistence.NewPlanRepository(db)))
	pinHandler := handler.NewPinHandler(pinuc.New(persistence.NewPinRepository(db)))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	r.Post("/plans", planHandler.Create)
	r.Get("/plans/{share_token}", planHandler.GetByShareToken)

	r.Post("/plans/{plan_id}/pins", pinHandler.Create)
	r.Patch("/plans/{plan_id}/pins/{pin_id}", pinHandler.Update)
	r.Delete("/plans/{plan_id}/pins/{pin_id}", pinHandler.Delete)

	return r
}
