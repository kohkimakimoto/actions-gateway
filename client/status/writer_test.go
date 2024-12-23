package status

import (
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_Init(t *testing.T) {
	d := testTempDir(t)
	f := filepath.Join(d, "status.json")
	w := NewWriter(f)
	err := w.Init()
	assert.NoError(t, err)

	// check the file exists
	_, err = os.Stat(f)
	assert.NoError(t, err)

	// check the status
	r := NewReader(f)
	s, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, CodeInactive, s.StatusCode)
}

func TestWriter_UpdateToConnecting(t *testing.T) {
	d := testTempDir(t)
	f := filepath.Join(d, "status.json")
	w := NewWriter(f)
	err := w.Init()
	assert.NoError(t, err)

	err = w.UpdateToConnecting(&types.SessionNewRequest{})
	assert.NoError(t, err)

	// check the status
	r := NewReader(f)
	s, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, CodeConnecting, s.StatusCode)
}
