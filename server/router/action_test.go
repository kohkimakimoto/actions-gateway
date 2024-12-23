package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionMessageFactory_NewMessage(t *testing.T) {
	f := NewActionMessageFactory(WithActionId("00000000-0000-0000-0000-000000000001"))
	msg, err := f.NewMessage("action name", "action body")
	assert.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", msg.Id)
	assert.Equal(t, "action name", msg.Name)
	assert.Equal(t, "action body", msg.Body)
}
