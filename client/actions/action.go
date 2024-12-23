package actions

import (
	"bytes"
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/client/config"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Action represents an action that is an executable file to be run by the client
type Action struct {
	Name string
	Path string
}

// ActionManager is a object that manages actions
type ActionManager struct {
	config    *config.Config
	actions   []*Action
	actionMap map[string]*Action
	spec      string
}

// NewActionManager creates a new ActionManager instance
func NewActionManager(cfg *config.Config) (*ActionManager, error) {
	// load actions
	actions, err := loadActions(cfg.ActionsAbsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load actions: %w", err)
	}
	// construct the action map
	actionMap := make(map[string]*Action)
	for _, a := range actions {
		actionMap[a.Name] = a
	}

	return &ActionManager{
		config:    cfg,
		actions:   actions,
		actionMap: actionMap,
	}, nil
}

// GetAction returns an action by Name
func (m *ActionManager) GetAction(name string) *Action {
	return m.actionMap[name]
}

func (m *ActionManager) Actions() []*Action {
	return m.actions
}

func (m *ActionManager) ActionNames() []string {
	names := make([]string, 0, len(m.actions))
	for _, a := range m.actions {
		names = append(names, a.Name)
	}
	return names
}

func (m *ActionManager) OutputSpec(errWriter io.Writer) (string, error) {
	if m.spec != "" {
		return m.spec, nil
	}

	actionPathSpecs := make([]*ActionPathSpec, 0)
	for _, a := range m.actions {
		spec, err := NewActionRunner(a, m.config.Dir(), errWriter).PathSpec()
		if err != nil {
			return "", fmt.Errorf("failed to generate spec for action %s: %w", a.Name, err)
		}
		if spec.Spec != "" {
			actionPathSpecs = append(actionPathSpecs, spec)
		}
	}

	funcMap := template.FuncMap{
		"indent": indent,
	}

	t := template.Must(template.New("spec").Funcs(funcMap).Parse(specTemplate))

	data := map[string]interface{}{
		"ActionPathSpecs": actionPathSpecs,
		"Config":          m.config,
	}

	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to generate spec from template: %w", err)
	}
	m.spec = buf.String()

	return m.spec, nil
}

// isExecutable checks if the file is executable
func isExecutable(mode os.FileMode) bool {
	return mode&0111 != 0 // Checks if any execution permission is set (user, group, others)
}

// loadActions retrieves the list of actions from the directory
func loadActions(dir string) ([]*Action, error) {
	var actions []*Action

	// Walks through the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If it's a file and is executable, register it as an action
		if !info.IsDir() && isExecutable(info.Mode()) {
			actions = append(actions, &Action{
				Name: info.Name(),
				Path: path,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return actions, nil
}

func indent(indent string, text string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if line != "" {
			result.WriteString(indent + line + "\n")
		}
	}
	return result.String()
}

// Note:
// "schemas: {}" is required to avoid the following error:
// https://community.openai.com/t/about-in-components-section-schemas-subsection-is-not-an-object/615947

var specTemplate = strings.TrimLeft(`
openapi: 3.1.0
info:
  title: {{ .Config.SpecInfo.Title }}
  summary: {{ .Config.SpecInfo.Summary }}
  description: {{ .Config.SpecInfo.Description }}
  version: {{ .Config.SpecInfo.Version }}
servers:
  - url: {{ .Config.ServerApiURL }}
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas: {}
security:
  - bearerAuth: []
paths:
{{- range .ActionPathSpecs }}
{{ indent "  " .ApiPath -}}
{{ indent "    " "post:" -}}
{{ indent "      " .Spec -}}
{{- end -}}
`, "\n")
