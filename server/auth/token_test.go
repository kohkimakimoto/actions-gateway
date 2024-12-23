package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenGenerator_NewTokenAsJWTString(t *testing.T) {
	k, _ := LoadKeyString("12345678901234567890123456789012")
	g := NewTokenGenerator(k, WithTime(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)), WithClientId("00000000-0000-0000-0000-000000000001"))
	token, err := g.NewTokenAsJWTString()
	assert.NoError(t, err)
	// t.Logf("Generated token: %s", token)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2MDk0NTkyMDAsInN1YiI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMSJ9.33F9jqdaYWuJWiz3w68f7SwhVHZ3fqkqgQ1ofQtQ1bY", token)
}
