package status

import (
	"encoding/json"
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"os"
	"sync"
)

type Writer struct {
	path   string
	status *Status
	mu     sync.RWMutex
}

func NewWriter(fPath string) *Writer {
	return &Writer{
		path:   fPath,
		status: initialStatus,
	}
}

func (m *Writer) Init() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.save()
}

func (m *Writer) UpdateToConnecting(req *types.SessionNewRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status.StatusCode = CodeConnecting
	m.status.SessionNewRequest = req
	m.status.SessionNewResponse = nil
	m.status.Error = ""

	return m.save()
}

func (m *Writer) UpdateToActive(req *types.SessionNewRequest, res *types.SessionNewResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status.StatusCode = CodeActive
	m.status.SessionNewRequest = req
	m.status.SessionNewResponse = res
	m.status.Error = ""

	return m.save()
}

func (m *Writer) UpdateToInactive(err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status.StatusCode = CodeInactive
	m.status.SessionNewRequest = nil
	m.status.SessionNewResponse = nil
	if err != nil {
		m.status.Error = err.Error()
	}

	return m.save()
}

// save saves the status to the file
// It is not thread-safe. You need to call this function with a lock.
func (m *Writer) save() error {
	if m.path == "" {
		return nil
	}
	b, err := json.Marshal(m.status)
	if err != nil {
		return fmt.Errorf("failed to serialize the status: %w", err)
	}
	if err := os.WriteFile(m.path, b, os.FileMode(0644)); err != nil {
		return fmt.Errorf("failed to write the status file: %w", err)
	}
	return nil
}
