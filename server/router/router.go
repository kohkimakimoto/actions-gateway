package router

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"net/http"
	"sync"
	"time"
)

// GenSessionIdFunc is a type of function to generate a session id
type GenSessionIdFunc func() (uuid.UUID, error)

// Router is an object to manage the sessions that are connected to the clients via websocket.
// A client is identified by the client id that is extracted from the JWT token provided by authentication.
// Each client can have only one session to transfer the actions to the proper client.
type Router struct {
	// sessions stores the sessions by client id
	sessions map[string]*Session
	// mutex for operations on sessions
	mu sync.RWMutex
	// timeout is an expiration time for inactive session
	timeout time.Duration
	// genSessionId is a function to generate a session id
	genSessionId GenSessionIdFunc
}

type Option func(*Router)

func WithSessionId(sId string) Option {
	return func(r *Router) {
		r.genSessionId = func() (uuid.UUID, error) {
			return uuid.Parse(sId)
		}
	}
}

// New creates a new Router object
func New(options ...Option) *Router {
	r := &Router{
		sessions:     make(map[string]*Session),
		timeout:      30 * time.Second,
		genSessionId: uuid.NewV7,
	}
	for _, option := range options {
		option(r)
	}
	return r
}

type SessionError struct {
	Message string
}

func NewSessionError(msg string) *SessionError {
	return &SessionError{
		Message: msg,
	}
}

func (e *SessionError) Error() string {
	return e.Message
}

var (
	ErrSessionAlreadyExists    = NewSessionError("Session already exists")
	ErrSessionAlreadyActivated = NewSessionError("Session already activated")
	ErrSessionNotFound         = NewSessionError("Session not found")
	ErrSessionInvalidId        = NewSessionError("Session id is invalid")
)

func (r *Router) NewSession(client *auth.Client, actions []string, spec string) (*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// check if the session already exists
	sess := r.sessions[client.Id]
	if sess != nil {
		return nil, ErrSessionAlreadyExists
	}

	// initialize a new session
	UUID, err := r.genSessionId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate a session id: %w", err)
	}
	sessionId := UUID.String()

	sess = &Session{}
	sess.id = sessionId
	sess.client = client
	sess.actions = actions
	sess.actionMap = make(map[string]bool)
	for _, a := range actions {
		sess.actionMap[a] = true
	}
	sess.spec = spec
	sess.results = make(map[string]chan *types.ActionResult)
	r.sessions[client.Id] = sess

	// start monitoring the session expiration
	// If the session is not activated within the expiration time, delete it
	go r.monitorInactiveSessionExpiration(sess)
	return sess, nil
}

func (r *Router) monitorInactiveSessionExpiration(sess *Session) {
	time.AfterFunc(r.timeout, func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		// if the session is not active, delete it
		if !sess.IsActive() {
			delete(r.sessions, sess.client.Id)
		}
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Router) ActivateSession(w http.ResponseWriter, req *http.Request, client *auth.Client, sessionId string) (*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// retrieve the session
	sess := r.sessions[client.Id]
	if sess == nil {
		return nil, ErrSessionNotFound
	}

	// check session id
	if sess.id != sessionId {
		return nil, ErrSessionInvalidId
	}

	// check if the session is already active
	if sess.IsActive() {
		return nil, ErrSessionAlreadyActivated
	}

	// upgrade http to websocket
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return nil, err
	}
	// save the websocket connection
	sess.conn = ws
	return sess, nil
}

func (r *Router) GetActiveSession(client *auth.Client) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sess := r.sessions[client.Id]
	if sess == nil || !sess.IsActive() {
		return nil
	}
	return sess
}

func (r *Router) CloseSession(sess *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sess.conn != nil {
		_ = sess.conn.Close()
	}

	delete(r.sessions, sess.client.Id)
}

func (r *Router) NumSessions() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions)
}
