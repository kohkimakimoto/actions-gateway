package router

import (
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession_Key(t *testing.T) {
	sess := &Session{
		id: "00000000-0000-0000-0000-000000000002",
		client: &auth.Client{
			Id: "00000000-0000-0000-0000-000000000001",
		},
	}
	assert.Equal(t, "00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002", sess.Key())
}

func TestSession_ConnectPath(t *testing.T) {
	sess := &Session{
		id: "00000000-0000-0000-0000-000000000002",
		client: &auth.Client{
			Id: "00000000-0000-0000-0000-000000000001",
		},
	}
	assert.Equal(t, "/api/session/connect/00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002", sess.ConnectPath())
}

func TestSession_IsActionExist(t *testing.T) {
	sess := &Session{
		actionMap: map[string]bool{
			"action1": true,
			"action2": true,
		},
	}
	assert.True(t, sess.IsActionExist("action1"))
	assert.True(t, sess.IsActionExist("action2"))
	assert.False(t, sess.IsActionExist("action3"))
}

func TestSession_AllocateResultChannel(t *testing.T) {
	// test allocate result channel
	sess := &Session{
		results: make(map[string]chan *types.ActionResult),
	}
	msgId := "00000000-0000-0000-0000-000000000000"
	ch := sess.AllocateResultChannel(msgId)
	assert.NotNil(t, ch)
	assert.NotNil(t, sess.results[msgId])
	assert.Equal(t, 1, len(sess.results))

	// test handle action result
	result := &types.ActionResult{
		Id:     msgId,
		Status: types.ActionResultStatusSuccess,
	}

	err := sess.HandleActionResult(result)
	assert.NoError(t, err)
	result2 := <-ch
	assert.Equal(t, result, result2)

	// test free result channel
	sess.FreeResultChannel(msgId)
	assert.Nil(t, sess.results[msgId])
	assert.Equal(t, 0, len(sess.results))

	// test handle action result that is already freed
	err = sess.HandleActionResult(result)
	assert.Error(t, err)
}
