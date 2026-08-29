package router

import (
	"database/sql"

	"github.com/Application-drop-up/Travellle/internal/handler"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/external"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/notification"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	noteuc "github.com/Application-drop-up/Travellle/internal/usecase/note"
	pinuc "github.com/Application-drop-up/Travellle/internal/usecase/pin"
	planuc "github.com/Application-drop-up/Travellle/internal/usecase/plan"
	spotuc "github.com/Application-drop-up/Travellle/internal/usecase/spot"
	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(db *sql.DB, googlePlacesAPIKey string, allowedOrigins []string, isDev bool) *chi.Mux {
	pinRepo := persistence.NewPinRepository(db)
	noteRepo := persistence.NewNoteRepository(db)

	pinUseCase := pinuc.New(pinRepo)
	noteUseCase := noteuc.New(noteRepo)

	planHandler := handler.NewPlanHandler(planuc.New(persistence.NewPlanRepository(db)), pinUseCase, noteUseCase)
	pinHandler := handler.NewPinHandler(pinUseCase)
	noteHandler := handler.NewNoteHandler(noteUseCase)
	spotHandler := handler.NewSpotHandler(spotuc.New(external.NewGooglePlacesClient(googlePlacesAPIKey), persistence.NewSpotRepository(db)))
	authHandler := handler.NewAuthHandler(useruc.New(
		persistence.NewUserRepository(db),
		persistence.NewSessionRepository(db),
		persistence.NewLoginOTPRepository(db),
		notification.NewLogEmailSender(),
	), isDev)

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
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

	mux.Route("/api/v1", func(r chi.Router) {
		r.Post("/user/register", authHandler.Register)
		r.Get("/user/{id}", authHandler.GetByID)
		r.Patch("/user/{id}", authHandler.Update)
		r.Delete("/user/{id}", authHandler.Delete)

		r.Post("/login", authHandler.LoginStart)
		r.Post("/login/verify", authHandler.LoginVerify)
		r.Post("/logout", authHandler.Logout)
		r.Get("/user/me", authHandler.Me)
	})

	return mux
}
