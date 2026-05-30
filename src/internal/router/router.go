package router

import (
	"database/sql"

	"github.com/Application-drop-up/Travellle/internal/handler"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/external"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	noteuc "github.com/Application-drop-up/Travellle/internal/usecase/note"
	pinuc "github.com/Application-drop-up/Travellle/internal/usecase/pin"
	planuc "github.com/Application-drop-up/Travellle/internal/usecase/plan"
	spotuc "github.com/Application-drop-up/Travellle/internal/usecase/spot"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB, googlePlacesAPIKey string) *chi.Mux {
	pinRepo := persistence.NewPinRepository(db)
	noteRepo := persistence.NewNoteRepository(db)

	pinUC := pinuc.New(pinRepo)
	noteUC := noteuc.New(noteRepo)

	planHandler := handler.NewPlanHandler(planuc.New(persistence.NewPlanRepository(db)), pinUC, noteUC)
	pinHandler := handler.NewPinHandler(pinUC)
	noteHandler := handler.NewNoteHandler(noteUC)
	spotHandler := handler.NewSpotHandler(spotuc.New(external.NewGooglePlacesClient(googlePlacesAPIKey)))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	r.Get("/spots/search", spotHandler.Search)

	r.Post("/plans", planHandler.Create)
	r.Get("/plans/{share_token}", planHandler.GetByShareToken)

	r.Get("/plans/{plan_id}/pins", pinHandler.List)
	r.Post("/plans/{plan_id}/pins", pinHandler.Create)
	r.Patch("/plans/{plan_id}/pins/{pin_id}", pinHandler.Update)
	r.Delete("/plans/{plan_id}/pins/{pin_id}", pinHandler.Delete)

	r.Post("/plans/{plan_id}/pins/{pin_id}/notes", noteHandler.Create)
	r.Patch("/plans/{plan_id}/pins/{pin_id}/notes/{note_id}", noteHandler.Update)
	r.Delete("/plans/{plan_id}/pins/{pin_id}/notes/{note_id}", noteHandler.Delete)

	return r
}
