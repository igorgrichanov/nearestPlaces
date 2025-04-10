package usecase

import (
	"errors"
	"nearestPlaces/internal/entity"
)

var ErrInternal = errors.New("internal server error")

type Auther interface {
	GetToken() (string, error)
}

type Storer interface {
	CreateIndexWithMapping() error
	UploadPlaces() error
}

type Restaurateur interface {
	GetPage(pageNum int) (*PageInfoDTO, error)
	GetClosestRestaurants(lat, lon float64) (*PageInfoDTO, error)
}

type PageInfoDTO struct {
	Name     string               `json:"name"`
	Total    int                  `json:"total,omitempty"`
	Places   []*entity.Restaurant `json:"places"`
	Page     int                  `json:"-"`
	PrevPage int                  `json:"prev_page,omitempty"`
	NextPage int                  `json:"next_page,omitempty"`
	LastPage int                  `json:"last_page,omitempty"`
}
