package main

import (
	"nearestPlaces/internal/app"
	"nearestPlaces/internal/lib/config"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
