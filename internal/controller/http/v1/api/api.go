package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"html/template"
	"log/slog"
	"nearestPlaces/internal/lib/api/response"
	"nearestPlaces/internal/lib/logger/sl"
	"nearestPlaces/internal/usecase"
	"net/http"
	"strconv"
)

type APIer interface {
	Places(w http.ResponseWriter, r *http.Request)
	Recommend(w http.ResponseWriter, r *http.Request)
	Paginate(http.ResponseWriter, *http.Request)
}

type Controller struct {
	log *slog.Logger
	uc  usecase.Restaurateur
}

func New(log *slog.Logger, uc usecase.Restaurateur) *Controller {
	return &Controller{
		log: log,
		uc:  uc,
	}
}

func (c *Controller) Places(w http.ResponseWriter, r *http.Request) {
	const op = "controller.places.Places"
	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		log.Error("invalid page number", slog.String("page", r.URL.Query().Get("page")))
		resp := fmt.Sprintf("Invalid 'page' value: '%d'.", page)
		render.Render(w, r, response.ErrBadRequest(resp))
		return
	}
	log.Info("request received", slog.String("page", r.URL.Query().Get("page")))

	pageInfo, err := c.uc.GetPage(page)
	if err != nil {
		log.Error("failed to get places: ", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
	log.Info("pages received from storage")

	if page > pageInfo.LastPage {
		log.Error("page parameter is too large", slog.String("page", r.URL.Query().Get("page")))
		resp := fmt.Sprintf("Invalid 'page' value: '%d'.", page)
		render.Render(w, r, response.ErrBadRequest(resp))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(pageInfo)
	if err != nil {
		log.Error("failed to encode response: ", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
}

func (c *Controller) Paginate(w http.ResponseWriter, r *http.Request) {
	const op = "controller.root.paginate"
	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	p, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || p < 1 {
		log.Error("invalid page number", slog.String("page", r.URL.Query().Get("page")))
		w.WriteHeader(http.StatusBadRequest)
		resp := fmt.Sprintf("Invalid 'page' value: '%d'.", p)
		w.Write([]byte(resp))
		return
	}
	log.Info("request received", slog.String("page", r.URL.Query().Get("page")))

	pageInfo, err := c.uc.GetPage(p)
	if err != nil {
		log.Error("failed to get places: ", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
	log.Info("pages received from storage")

	if p > pageInfo.LastPage {
		log.Error("page parameter is too large", slog.String("page", r.URL.Query().Get("page")))
		w.WriteHeader(http.StatusBadRequest)
		resp := fmt.Sprintf("Invalid 'page' value: '%d'.", p)
		w.Write([]byte(resp))
		return
	}
	tmpl, err := template.New("index.html").ParseFiles("templates/index.html")
	if err != nil {
		log.Error("failed to parse template: ", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}

	err = tmpl.Execute(w, pageInfo)
	if err != nil {
		log.Error("failed to render template: ", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
	log.Info("request executed")
}

func (c *Controller) Recommend(w http.ResponseWriter, r *http.Request) {
	const op = "controller.recommend.Recommend"
	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil || lat < 0 {
		log.Error("failed to parse latitude")
		resp := fmt.Sprintf("Invalid 'lat' value: '%v'.", r.URL.Query().Get("lat"))
		render.Render(w, r, response.ErrBadRequest(resp))
		return
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil || lon < 0 {
		log.Error("failed to parse longitude")
		resp := fmt.Sprintf("Invalid 'lon' value: '%v'.", r.URL.Query().Get("lon"))
		render.Render(w, r, response.ErrBadRequest(resp))
		return
	}
	result, err := c.uc.GetClosestRestaurants(lat, lon)
	if err != nil {
		log.Error("failed to get closest restaurants", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error("failed to encode response: ", sl.Err(err))
		render.Render(w, r, response.ErrInternal())
		return
	}
}
