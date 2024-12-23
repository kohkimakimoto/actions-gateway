package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadKey(t *testing.T) {
	// Test case: secret is less than 32 bytes
	k, err := LoadKey([]byte("1234567890"))
	assert.Error(t, err)
	assert.Nil(t, k)

	// Test case: secret is 32 bytes
	k, err = LoadKey([]byte("12345678901234567890123456789012"))
	assert.NoError(t, err)
	assert.NotNil(t, k)
}

func TestLoadKeyString(t *testing.T) {
	// Test case: secret is less than 32 bytes
	k, err := LoadKeyString("1234567890")
	assert.Error(t, err)
	assert.Nil(t, k)

	// Test case: secret is 32 bytes
	k, err = LoadKeyString("12345678901234567890123456789012")
	assert.NoError(t, err)
	assert.NotNil(t, k)
}
