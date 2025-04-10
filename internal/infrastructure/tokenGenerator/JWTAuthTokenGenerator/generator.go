package JWTAuthTokenGenerator

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/jwtauth/v5"
	"nearestPlaces/internal/infrastructure/tokenGenerator"
	"time"
)

type JWTAuth struct {
	TokenAuth     *jwtauth.JWTAuth
	tokenLiveTime time.Duration
}

func New(tokenAuth *jwtauth.JWTAuth, tokenLiveTime time.Duration) *JWTAuth {
	return &JWTAuth{
		TokenAuth:     tokenAuth,
		tokenLiveTime: tokenLiveTime,
	}
}

func (m *JWTAuth) Generate() (string, error) {
	_, tokenString, err := m.TokenAuth.Encode(map[string]interface{}{
		"iss": "localhost:8888",
		//"sub": userLogin,
		"aud": "localhost:8888",
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(m.tokenLiveTime).Unix(),
		"jti": gofakeit.UUID(),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", tokenGenerator.GenerationError, err)
	}
	return tokenString, nil
}
