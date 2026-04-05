package router

import (
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/config"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/handlers"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
)

// NewRouter wires the chi router with all middlewares and routes.
func NewRouter(cfg *config.Config, db *pgxpool.Pool, logger *zap.Logger, authSvc service.AuthService, learning *handlers.LearningHandler, autopilot *handlers.AutopilotHandler, aiHandler *handlers.AIHandler, planner *handlers.PlannerHandler, tutoring *handlers.TutoringHandler, feed *handlers.FeedHandler, profile *handlers.ProfileHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(chimiddleware.Recoverer)

	health := &handlers.HealthHandler{DB: db}
	auth := &handlers.AuthHandler{Service: authSvc}

	r.Get("/healthz", health.Healthz)
	r.Get("/readyz", health.Readyz)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/auth/register", auth.Register)
		v1.Post("/auth/login", auth.Login)
		v1.Post("/auth/logout", auth.Logout)

		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(cfg.JWTSecret))
			protected.Get("/me", auth.Me)
			protected.Get("/profile/summary", profile.Summary)
			protected.Get("/autopilot/today", autopilot.Today)

			protected.Route("/courses", func(cr chi.Router) {
				cr.Get("/", learning.ListCourses)
				cr.Post("/", learning.CreateCourse)
				cr.Post("/import-from-lms", learning.ImportFromLMS)
				cr.Get("/{courseId}", learning.GetCourse)
				cr.Delete("/{courseId}", learning.DeleteCourse)
				cr.Get("/{courseId}/notes", learning.ListNotes)
				cr.Get("/{courseId}/notes/metrics", learning.ListNoteMetrics)
				cr.Get("/{courseId}/topics", learning.ListTopics)
				cr.Post("/{courseId}/topics", learning.CreateTopic)
				cr.Get("/{courseId}/flashcards", learning.ListFlashcards)
			})
			protected.Route("/notes", func(nr chi.Router) {
				nr.Post("/", learning.CreateNote)
				nr.Post("/{noteId}/ask", learning.AskNoteAI)
				nr.Get("/{noteId}", learning.GetNote)
				nr.Put("/{noteId}", learning.UpdateNote)
				nr.Delete("/{noteId}", learning.DeleteNote)
			})
			protected.Post("/flashcards", learning.CreateFlashcard)
			protected.Delete("/flashcards/{cardId}", learning.DeleteFlashcard)
			protected.Get("/review/due", learning.GetDueCards)
			protected.Post("/review/events", learning.SubmitReview)
			protected.Get("/insights/confusion", learning.GetConfusionInsights)
			protected.Post("/flashcards/generate", aiHandler.GenerateFlashcards)
			protected.Post("/flashcards/generate/accept", aiHandler.AcceptGeneratedFlashcards)

			protected.Route("/planner", func(pr chi.Router) {
				pr.Get("/sessions", planner.ListSessions)
				pr.Post("/sessions", planner.CreateSession)
				pr.Patch("/sessions/{id}/status", planner.UpdateSessionStatus)
				pr.Delete("/sessions/{id}", planner.DeleteSession)
				pr.Post("/missed-yesterday", planner.MarkMissed)
				pr.Get("/deadlines", planner.ListDeadlines)
				pr.Post("/deadlines", planner.CreateDeadline)
				pr.Delete("/deadlines/{id}", planner.DeleteDeadline)
				pr.Post("/regenerate", planner.RegeneratePlan)
			})

			protected.Route("/tutoring", func(tr chi.Router) {
				tr.Get("/teaching-topics", tutoring.ListTeachingTopics)
				tr.Get("/requests/incoming", tutoring.ListIncomingRequests)
				tr.Get("/requests/outgoing", tutoring.ListOutgoingRequests)
				tr.Post("/requests", tutoring.CreateRequest)
				tr.Patch("/requests/{id}", tutoring.UpdateRequestStatus)
				tr.Get("/match", tutoring.FindTutorMatches)
			})

			protected.Route("/feed", func(fr chi.Router) {
				fr.Get("/", feed.ListFeedPosts)
				fr.Post("/", feed.CreateFeedPost)
				fr.Get("/{postId}/comments", feed.ListFeedComments)
				fr.Post("/{postId}/comments", feed.CreateFeedComment)
				fr.Delete("/{postId}", feed.DeleteFeedPost)
				fr.Post("/{postId}/upvote", feed.ToggleFeedUpvote)
			})
		})
	})

	return r
}
