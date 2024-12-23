package auth

import (
	"errors"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// LoadKeyString creates a new key from the given secret.
func LoadKeyString(secret string) (jwk.Key, error) {
	return LoadKey([]byte(secret))
}

// LoadKey creates a new key from the given secret.
func LoadKey(secret []byte) (jwk.Key, error) {
	if len(secret) < 32 {
		return nil, errors.New("secret must be at least 32 bytes")
	}
	return jwk.FromRaw(secret)
}
