package auth

import (
	"errors"
	"log/slog"
	"nearestPlaces/internal/infrastructure/tokenGenerator"
	"nearestPlaces/internal/lib/logger/sl"
)

var ErrInternal = errors.New("internal server error")

//go:generate go run github.com/vektra/mockery/v2@v2.53.0 --name=TokenGenerator
type TokenGenerator interface {
	Generate() (string, error)
}

type UseCase struct {
	log *slog.Logger
	tg  TokenGenerator
}

func New(log *slog.Logger, tg TokenGenerator) *UseCase {
	return &UseCase{
		log: log,
		tg:  tg,
	}
}

func (u *UseCase) GetToken() (string, error) {
	const op = "service.auth.Login"
	log := u.log.With(
		slog.String("op", op),
	)
	t, err := u.tg.Generate()
	if errors.Is(err, tokenGenerator.GenerationError) {
		log.Error("error generating token", sl.Err(err))
		return "", ErrInternal
	} else if err != nil {
		log.Error("unable to generate token", sl.Err(err))
		return "", ErrInternal
	}
	log.Info("token generated", sl.Info(t))
	return t, nil
}
