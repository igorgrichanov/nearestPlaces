package controller

import (
	apiController "nearestPlaces/internal/controller/http/v1/api"
	authController "nearestPlaces/internal/controller/http/v1/auth"
)

type Controllers struct {
	Auth authController.Auther
	Api  apiController.APIer
}

func New(auth authController.Auther, api apiController.APIer) *Controllers {
	return &Controllers{
		Auth: auth,
		Api:  api,
	}
}
