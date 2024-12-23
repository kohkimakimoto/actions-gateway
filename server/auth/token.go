package auth

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

type Token struct {
	// Subject (JWT claim "sub")
	// This value is used as the "Client Identifier" in this system.
	Subject string
	// IssuedAt (JWT claim "iat")
	IssuedAt time.Time
}

type GenClientIdFunc func() (uuid.UUID, error)

type TokenGenerator struct {
	key         jwk.Key
	genClientId GenClientIdFunc
	clock       func() time.Time
}

type TokenGeneratorOption func(*TokenGenerator)

func WithClientId(cId string) TokenGeneratorOption {
	return func(g *TokenGenerator) {
		g.genClientId = func() (uuid.UUID, error) {
			return uuid.Parse(cId)
		}
	}
}

func WithTime(t time.Time) TokenGeneratorOption {
	return func(g *TokenGenerator) {
		g.clock = func() time.Time {
			return t
		}
	}
}

func NewTokenGenerator(key jwk.Key, options ...TokenGeneratorOption) *TokenGenerator {
	g := &TokenGenerator{
		key:         key,
		genClientId: uuid.NewV7,
		clock:       time.Now,
	}
	for _, option := range options {
		option(g)
	}
	return g
}

func (g *TokenGenerator) NewTokenAsJWTString() (string, error) {
	UUID, err := g.genClientId()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}
	clientId := UUID.String()

	t := &Token{
		Subject:  clientId,
		IssuedAt: g.clock(),
	}

	jwtToken, err := jwt.NewBuilder().
		Subject(t.Subject).
		IssuedAt(t.IssuedAt).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build JWT token: %w", err)
	}

	// Sign the token using HS256 algorithm
	signedToken, err := jwt.Sign(jwtToken, jwt.WithKey(jwa.HS256, g.key))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return string(signedToken), nil
}
