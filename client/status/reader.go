package status

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
)

type Reader struct {
	path string
}

func NewReader(fPath string) *Reader {
	return &Reader{
		path: fPath,
	}
}

func (r *Reader) Watcher() (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create a fsnotify.Reader: %w", err)
	}
	if err := w.Add(r.path); err != nil {
		return nil, fmt.Errorf("failed to watch the status file: %w", err)
	}
	return w, nil
}

func (r *Reader) Read() (*Status, error) {
	status := &Status{}
	b, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return initialStatus, nil
		}
		return nil, fmt.Errorf("failed to read status file: %w", err)
	}
	if err := json.Unmarshal(b, status); err != nil {
		return nil, fmt.Errorf("failed to deserialize the status: %w", err)
	}
	return status, nil
}
