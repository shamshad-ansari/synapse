package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/mock-canvas/internal/handlers"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8082"
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Sync() //nolint:errcheck

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"mock-canvas"}`))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(handlers.RequireToken)

		api.Get("/users/self", handlers.GetSelf)
		api.Get("/courses", handlers.ListCourses)
		api.Get("/courses/{courseID}", handlers.GetCourse)
		api.Get("/courses/{courseID}/assignments", handlers.ListAssignments)
		api.Get("/courses/{courseID}/students/submissions", handlers.ListSubmissions)
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("mock-canvas starting on :" + port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
	logger.Info("mock-canvas stopped")
}
