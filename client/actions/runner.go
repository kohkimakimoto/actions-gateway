package actions

import (
	"bytes"
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// ActionRunner runs an action
type ActionRunner struct {
	action    *Action
	workDir   string
	errWriter io.Writer
}

func NewActionRunner(action *Action, workDir string, errWriter io.Writer) *ActionRunner {
	return &ActionRunner{
		action:    action,
		workDir:   workDir,
		errWriter: errWriter,
	}
}

type ActionPathSpec struct {
	Name    string
	ApiPath string
	Spec    string
}

// PathSpec returns the OpenAPI Path spec of this action's endpoint
func (r *ActionRunner) PathSpec() (*ActionPathSpec, error) {
	ex, err := r.resolveExecutablePath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(r.action.Path)
	cmd.Dir = r.workDir
	cmd.Stderr = r.errWriter
	cmd.Env = append(os.Environ(), "ACTIONS_GATEWAY_EXECUTABLE="+ex, "ACTIONS_GATEWAY_ACTIONS_SPEC=1")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to generate spec: %w", err)
	}

	return &ActionPathSpec{
		Name:    r.action.Name,
		ApiPath: path.Join("/actions", r.action.Name+":"),
		Spec:    string(output),
	}, nil
}

// Run runs the action
func (r *ActionRunner) Run(msg *types.ActionMessage) ([]byte, error) {
	ex, err := r.resolveExecutablePath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(r.action.Path)
	cmd.Dir = r.workDir
	cmd.Stderr = r.errWriter
	cmd.Env = append(os.Environ(), "ACTIONS_GATEWAY_EXECUTABLE="+ex)
	cmd.Stdin = bytes.NewReader([]byte(msg.Body))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run action: %w", err)
	}
	return output, nil
}

// resolveExecutablePath resolves the Path of the executable
func (r *ActionRunner) resolveExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable Path: %w", err)
	}
	ex, err = filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlink: %w", err)
	}
	return ex, nil
}
