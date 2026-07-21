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
	"github.com/go-chi/cors"
)

func New(db *sql.DB, googlePlacesAPIKey string, allowedOrigins []string) *chi.Mux {
	pinRepo := persistence.NewPinRepository(db)
	noteRepo := persistence.NewNoteRepository(db)

	pinUC := pinuc.New(pinRepo)
	noteUC := noteuc.New(noteRepo)

	planHandler := handler.NewPlanHandler(planuc.New(persistence.NewPlanRepository(db)), pinUC, noteUC)
	pinHandler := handler.NewPinHandler(pinUC)
	noteHandler := handler.NewNoteHandler(noteUC)
	spotHandler := handler.NewSpotHandler(spotuc.New(external.NewGooglePlacesClient(googlePlacesAPIKey)))

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/health", handler.Health)

	mux.Get("/spots/search", spotHandler.Search)

	mux.Post("/plans", planHandler.Create)
	mux.Get("/plans/{share_token}", planHandler.GetByShareToken)

	mux.Get("/plans/{plan_id}/pins", pinHandler.List)
	mux.Post("/plans/{plan_id}/pins", pinHandler.Create)
	mux.Patch("/plans/{plan_id}/pins/{pin_id}", pinHandler.Update)
	mux.Delete("/plans/{plan_id}/pins/{pin_id}", pinHandler.Delete)

	mux.Post("/plans/{plan_id}/pins/{pin_id}/notes", noteHandler.Create)
	mux.Patch("/plans/{plan_id}/pins/{pin_id}/notes/{note_id}", noteHandler.Update)
	mux.Delete("/plans/{plan_id}/pins/{pin_id}/notes/{note_id}", noteHandler.Delete)

	return mux
}
