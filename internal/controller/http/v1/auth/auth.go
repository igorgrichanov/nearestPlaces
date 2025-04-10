package auth

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"nearestPlaces/internal/lib/api/response"
	"nearestPlaces/internal/lib/logger/sl"
	"nearestPlaces/internal/usecase"
	"net/http"
)

type Auther interface {
	GetToken(http.ResponseWriter, *http.Request)
}

type Controller struct {
	log *slog.Logger
	uc  usecase.Auther
}

func New(log *slog.Logger, uc usecase.Auther) *Controller {
	return &Controller{
		log: log,
		uc:  uc,
	}
}

type Response struct {
	Token string `json:"token"`
}

func (c *Controller) GetToken(w http.ResponseWriter, r *http.Request) {
	const op = "controller.auth.GetToken"
	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	log.Info("request received")
	token, err := c.uc.GetToken()
	if err != nil {
		log.Error("failed to get token", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
	log.Info("token generated")

	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Token: token,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
}
