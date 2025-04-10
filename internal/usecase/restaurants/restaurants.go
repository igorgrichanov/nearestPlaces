package restaurants

import (
	"log/slog"
	"nearestPlaces/internal/entity"
	"nearestPlaces/internal/lib/logger/sl"
	"nearestPlaces/internal/usecase"
)

type UseCase struct {
	log     *slog.Logger
	storage Store
}

func New(log *slog.Logger, storage Store) *UseCase {
	return &UseCase{
		log:     log,
		storage: storage,
	}
}

type Store interface {
	GetClosest(lat, lon float64) ([]*entity.Restaurant, error)
	GetPlaces(limit, offset int) ([]*entity.Restaurant, int, error)
}

func (u *UseCase) GetClosestRestaurants(lat, lon float64) (*usecase.PageInfoDTO, error) {
	const op = "usecase.restaurants.GetClosestRestaurants"
	log := u.log.With(
		slog.String("op", op),
	)
	places, err := u.storage.GetClosest(lat, lon)
	if err != nil {
		log.Error("failed to get closest restaurants", sl.Err(err))
		return nil, usecase.ErrInternal
	}
	log.Info("closest restaurants received")

	result := &usecase.PageInfoDTO{
		Name:   "Recommendation",
		Places: places,
	}
	return result, nil
}

func (u *UseCase) GetPage(pageNum int) (*usecase.PageInfoDTO, error) {
	const op = "usecase.restaurants.GetPages"
	log := u.log.With(
		slog.String("op", op),
	)
	limit := 10
	offset := (pageNum - 1) * limit
	places, total, err := u.storage.GetPlaces(limit, offset)
	if err != nil {
		log.Error("failed to get places: ", sl.Err(err))
		return nil, usecase.ErrInternal
	}
	log.Info("page received from storage")

	result := &usecase.PageInfoDTO{
		Name:     "Places",
		Total:    total,
		Places:   places,
		Page:     pageNum,
		PrevPage: pageNum - 1,
		NextPage: pageNum + 1,
		LastPage: total/limit + 1,
	}
	return result, nil
}
