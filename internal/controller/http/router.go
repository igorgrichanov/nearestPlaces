package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"nearestPlaces/internal/controller"
	"nearestPlaces/internal/controller/http/middleware/logger"
	"net/http"
)

func NewRouter(log *slog.Logger, ctrl *controller.Controllers, ja *jwtauth.JWTAuth) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(logger.New(log))
	router.Get("/", ctrl.Api.Paginate)
	router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(ja))
			r.Use(jwtauth.Authenticator(ja))
			r.Get("/recommend", ctrl.Api.Recommend)
		})

		r.Get("/places", ctrl.Api.Places)
		r.Get("/get_token", ctrl.Auth.GetToken)
	})
	return router
}
