package actions

import (
	"github.com/kohkimakimoto/actions-gateway/client/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestActionManager_GetAction(t *testing.T) {
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

	m, err := NewActionManager(&config.Config{
		ActionsAbsDir: actionsDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	a := m.GetAction("testAction")
	assert.NotNil(t, a)
	assert.Equal(t, "testAction", a.Name)
}

func TestActionManager_Actions(t *testing.T) {
	dir := testTempDir(t)
	actionsDir := filepath.Join(dir, "actions")
	err := os.MkdirAll(actionsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	testAction1File := filepath.Join(actionsDir, "testAction1")
	err = os.WriteFile(testAction1File, []byte(`#!/usr/bin/env bash
echo 'This is a test action1'
`), 0755)
	if err != nil {
		t.Fatal(err)
	}
	testAction2File := filepath.Join(actionsDir, "testAction2")
	err = os.WriteFile(testAction2File, []byte(`#!/usr/bin/env bash
echo 'This is a test action2'
`), 0755)
	if err != nil {
		t.Fatal(err)
	}
	m, err := NewActionManager(&config.Config{
		ActionsAbsDir: actionsDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	as := m.Actions()
	assert.Len(t, as, 2)
	assert.Equal(t, "testAction1", as[0].Name)
	assert.Equal(t, "testAction2", as[1].Name)
}

func TestActionManager_ActionNames(t *testing.T) {
	dir := testTempDir(t)
	actionsDir := filepath.Join(dir, "actions")
	err := os.MkdirAll(actionsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	testAction1File := filepath.Join(actionsDir, "testAction1")
	err = os.WriteFile(testAction1File, []byte(`#!/usr/bin/env bash
echo 'This is a test action1'
`), 0755)
	if err != nil {
		t.Fatal(err)
	}
	testAction2File := filepath.Join(actionsDir, "testAction2")
	err = os.WriteFile(testAction2File, []byte(`#!/usr/bin/env bash
echo 'This is a test action2'
`), 0755)
	if err != nil {
		t.Fatal(err)
	}
	m, err := NewActionManager(&config.Config{
		ActionsAbsDir: actionsDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	names := m.ActionNames()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "testAction1")
	assert.Contains(t, names, "testAction2")
}

// TODO: Add more tests
