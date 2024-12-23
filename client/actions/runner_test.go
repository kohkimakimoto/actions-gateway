package actions

import (
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestActionRunner_PathSpec(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := testTempDir(t)
		actionsDir := filepath.Join(dir, "actions")
		err := os.MkdirAll(actionsDir, 0755)
		if err != nil {
			t.Fatal(err)
		}
		testActionFile := filepath.Join(actionsDir, "testAction")
		err = os.WriteFile(testActionFile, []byte(`#!/usr/bin/env bash
if [[ -n "$ACTIONS_GATEWAY_ACTIONS_SPEC" ]]; then
  echo 'You should output the OpenAPI spec of this action here'
fi
`), 0755)
		if err != nil {
			t.Fatal(err)
		}
		action := &Action{
			Name: "testAction",
			Path: testActionFile,
		}

		spec, err := NewActionRunner(action, dir, nil).PathSpec()
		assert.NoError(t, err)
		assert.Equal(t, "testAction", spec.Name)
		assert.Equal(t, "/actions/testAction:", spec.ApiPath)
		assert.Equal(t, "You should output the OpenAPI spec of this action here\n", spec.Spec)
	})
}

func TestActionRunner_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := testTempDir(t)
		actionsDir := filepath.Join(dir, "actions")
		err := os.MkdirAll(actionsDir, 0755)
		if err != nil {
			t.Fatal(err)
		}
		testActionFile := filepath.Join(actionsDir, "testAction")
		err = os.WriteFile(testActionFile, []byte(`#!/usr/bin/env bash
echo 'This is a test action'
`), 0755)
		if err != nil {
			t.Fatal(err)
		}
		action := &Action{
			Name: "testAction",
			Path: testActionFile,
		}

		b, err := NewActionRunner(action, dir, nil).Run(&types.ActionMessage{
			Id:   "00000000-0000-0000-0000-000000000001",
			Name: "testAction",
		})
		assert.NoError(t, err)
		assert.Equal(t, "This is a test action\n", string(b))
	})
}
