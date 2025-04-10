package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"log/slog"
	"nearestPlaces/internal/controller"
	httpController "nearestPlaces/internal/controller/http"
	"nearestPlaces/internal/controller/http/v1/api"
	authController "nearestPlaces/internal/controller/http/v1/auth"
	"nearestPlaces/internal/infrastructure/JSONSchemaReader"
	"nearestPlaces/internal/infrastructure/csv"
	"nearestPlaces/internal/infrastructure/repository/elastic"
	"nearestPlaces/internal/infrastructure/tokenGenerator/JWTAuthTokenGenerator"
	"nearestPlaces/internal/lib/config"
	"nearestPlaces/internal/lib/logger/sl"
	"nearestPlaces/internal/usecase/auth"
	"nearestPlaces/internal/usecase/restaurants"
	"nearestPlaces/internal/usecase/store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const indexName = "places"

func Run(cfg *config.Config) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("logger started successfully")

	// infrastructure

	// db
	elasticAddr := fmt.Sprintf("http://%s:%s", cfg.Elastic.Host, cfg.Elastic.Port)
	esConfig := elasticsearch.Config{
		Addresses: []string{elasticAddr},
	}
	es, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		log.Error("failed to create elasticsearch client: ", sl.Err(err))
		os.Exit(1)
	}

	// storage
	storage := elastic.New(log, es, indexName)

	mappingReader := JSONSchemaReader.New()
	csvParser := csv.New()

	ja := jwtauth.New("HS256", []byte(cfg.Token.Secret), nil,
		jwt.WithAcceptableSkew(cfg.Token.Skew))
	tokenGenerator := JWTAuthTokenGenerator.New(ja, cfg.Token.TTL)

	// use cases
	restaurantsUseCase := restaurants.New(log, storage)
	storeUseCase := store.New(log, cfg, mappingReader, csvParser, storage)
	authUseCase := auth.New(log, tokenGenerator)
	err = storeUseCase.CreateIndexWithMapping()
	if err != nil {
		log.Error("failed to create index: ", sl.Err(err))
	}
	err = storeUseCase.UploadPlaces()
	if err != nil {
		log.Error("failed to fill index: ", sl.Err(err))
	}

	// controller
	apiCtrl := api.New(log, restaurantsUseCase)
	authCtrl := authController.New(log, authUseCase)
	ctrl := controller.New(authCtrl, apiCtrl)

	// router
	router := httpController.NewRouter(log, ctrl, ja)

	// server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("starting server at http://" + addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err.Error())
		}
	}()

	<-done
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("error while shutting down server", sl.Err(err))
	}

	log.Info("shut down successfully")
}
