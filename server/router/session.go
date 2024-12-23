package router

import (
	"github.com/gorilla/websocket"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/pkg/errors"
	"sync"
)

// Session represents a connection with a client.
type Session struct {
	// id is a session id
	id string
	// client is a connected an authenticated client
	client *auth.Client
	// conn is a websocket connection
	conn *websocket.Conn
	// actions is a list of action names that are supported by the session
	actions []string
	// actionMap is a map of action names generated from the actions list.
	// It is used to retrieve the action name quickly.
	actionMap map[string]bool
	// spec is an OpenAPI Spec written in YAML.
	spec string
	// results is a map of result channels. The key of the map is an action message id (UUID v7).
	results map[string]chan *types.ActionResult
	// mu is a mutex for operations on results
	mu sync.RWMutex
}

func (sess *Session) Conn() *websocket.Conn {
	return sess.conn
}

func (sess *Session) Spec() string {
	return sess.spec
}

func (sess *Session) IsActive() bool {
	return sess.conn != nil
}

func (sess *Session) Key() string {
	return sess.client.Id + "/" + sess.id
}

func (sess *Session) ConnectPath() string {
	return "/api/session/connect/" + sess.Key()
}

func (sess *Session) IsActionExist(name string) bool {
	return sess.actionMap[name]
}

func (sess *Session) AllocateResultChannel(msgId string) <-chan *types.ActionResult {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	resultChan := make(chan *types.ActionResult, 1)
	sess.results[msgId] = resultChan

	return resultChan
}

func (sess *Session) FreeResultChannel(msgId string) {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if ch, ok := sess.results[msgId]; ok {
		close(ch)
		delete(sess.results, msgId)
	}
}

func (sess *Session) HandleActionResult(result *types.ActionResult) error {
	sess.mu.RLock()
	defer sess.mu.RUnlock()

	ch, ok := sess.results[result.Id]
	if !ok {
		return errors.New("result channel not found")
	}
	ch <- result
	return nil
}
