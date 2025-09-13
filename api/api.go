package api

import (
	"time"

	"github.com/LocTr/my-blog-be/api/app"
	"github.com/LocTr/my-blog-be/database"
	"github.com/LocTr/my-blog-be/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(enableCORS bool) (*chi.Mux, error) {
	logger := logging.NewLogger()

	logger.Info("Connecting to database...")

	db, err := database.DBConnect()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	appAPI, err := app.NewAPI(db)
	if err != nil {
		logger.WithField("module", "api").Error(err)
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(logging.NewStructuredLogger(logger)) // Use custom structured logger
	// router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(15 * time.Second))

	// use CORS if client is not served by this api
	if enableCORS {
		router.Use(corsConfig().Handler)
	}

	router.Group(func(router chi.Router) {
		router.Mount("/posts", appAPI.Router())
	})

	return router, nil
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}
