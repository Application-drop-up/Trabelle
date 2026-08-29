package router

import (
	"database/sql"

	"github.com/Application-drop-up/Travellle/internal/handler"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/external"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/notification"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/usecase/note"
	"github.com/Application-drop-up/Travellle/internal/usecase/pin"
	"github.com/Application-drop-up/Travellle/internal/usecase/plan"
	"github.com/Application-drop-up/Travellle/internal/usecase/spot"
	"github.com/Application-drop-up/Travellle/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(db *sql.DB, googlePlacesAPIKey string, allowedOrigins []string, isDev bool) *chi.Mux {
	pinRepo := persistence.NewPinRepository(db)
	noteRepo := persistence.NewNoteRepository(db)

	pinUseCase := pin.New(pinRepo)
	noteUseCase := note.New(noteRepo)

	planHandler := handler.NewPlanHandler(plan.New(persistence.NewPlanRepository(db)), pinUseCase, noteUseCase)
	pinHandler := handler.NewPinHandler(pinUseCase)
	noteHandler := handler.NewNoteHandler(noteUseCase)
	spotHandler := handler.NewSpotHandler(spot.New(external.NewGooglePlacesClient(googlePlacesAPIKey), persistence.NewSpotRepository(db)))
	authHandler := handler.NewAuthHandler(user.New(
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

	mux.Route("/api/v1", func(routing chi.Router) {
		routing.Post("/user/register", authHandler.Register)
		routing.Get("/user/{id}", authHandler.GetByID)
		routing.Patch("/user/{id}", authHandler.Update)
		routing.Delete("/user/{id}", authHandler.Delete)

		routing.Post("/login", authHandler.LoginStart)
		routing.Post("/login/verify", authHandler.LoginVerify)
		routing.Post("/logout", authHandler.Logout)
		routing.Get("/user/me", authHandler.Me)

		routing.Post("/plans/{share_token}/publish", planHandler.Publish)

		routing.Post("/user/{id}/spot/share", spotHandler.Save)
	})

	return mux
}
