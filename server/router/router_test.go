package router

import (
	"bufio"
	"bytes"
	"github.com/google/uuid"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRouter_NewSession(t *testing.T) {
	t.Run("new session", func(t *testing.T) {
		r := New()
		r.genSessionId = func() (uuid.UUID, error) {
			return uuid.MustParse("00000000-0000-0000-0000-000000000002"), nil
		}
		r.timeout = 1 * time.Second
		ct := &auth.Client{
			Id: "00000000-0000-0000-0000-000000000001",
		}
		sess, err := r.NewSession(ct, []string{"action1", "action2"}, "")
		assert.NoError(t, err)
		assert.NotNil(t, sess)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", sess.client.Id)
		assert.Equal(t, "00000000-0000-0000-0000-000000000002", sess.id)
		assert.Equal(t, 1, r.NumSessions())

		// test the session already exists
		sess, err = r.NewSession(ct, []string{"action1", "action2"}, "")
		assert.Error(t, err)
		assert.Nil(t, sess)
		assert.Equal(t, ErrSessionAlreadyExists, err)

		// test the session timeout
		time.Sleep(2 * time.Second)
		assert.Equal(t, 0, r.NumSessions()) // the expired session is removed

		// the session can be created again
		sess, err = r.NewSession(ct, []string{"action1", "action2"}, "")
		assert.NoError(t, err)
		assert.NotNil(t, sess)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", sess.client.Id)
		assert.Equal(t, "00000000-0000-0000-0000-000000000002", sess.id)
		assert.Equal(t, 1, r.NumSessions())
	})
}

func TestRouter_ActivateSession(t *testing.T) {
	t.Run("activate session", func(t *testing.T) {
		r := New()
		r.genSessionId = func() (uuid.UUID, error) {
			return uuid.MustParse("00000000-0000-0000-0000-000000000002"), nil
		}
		ct := &auth.Client{
			Id: "00000000-0000-0000-0000-000000000001",
		}

		// no active session
		assert.Nil(t, r.GetActiveSession(ct))

		// create a new session
		sess, err := r.NewSession(ct, []string{"action1", "action2"}, "")
		assert.NoError(t, err)

		// activate the session
		req := &http.Request{
			Method: http.MethodGet,
			Header: http.Header{
				"Upgrade":               []string{"websocket"},
				"Connection":            []string{"upgrade"},
				"Sec-Websocket-Key":     []string{"dGhlIHNhbXBsZSBub25jZQ=="},
				"Sec-Websocket-Version": []string{"13"},
			}}

		br := bufio.NewReaderSize(strings.NewReader(""), 4096)
		bw := bufio.NewWriterSize(&bytes.Buffer{}, 4096)
		resp := &testHijakableResponseWriter{
			brw: bufio.NewReadWriter(br, bw),
		}
		sess2, err := r.ActivateSession(resp, req, ct, sess.id)
		assert.NoError(t, err)
		assert.NotNil(t, sess2)
		assert.Equal(t, sess, sess2)
		assert.True(t, sess.IsActive())
		assert.NotNil(t, sess.Conn())

		// get the active session
		sess3 := r.GetActiveSession(ct)
		assert.NotNil(t, sess3)
		assert.Equal(t, sess, sess3)

		// close the session
		r.CloseSession(sess)
		assert.Nil(t, r.GetActiveSession(ct))
	})

}

// ----------------------------------------------------------------
// websocket test helpers
// The following code is referenced from github.com/gorilla/websocket package test code.
// ----------------------------------------------------------------

type fakeNetConn struct {
	io.Reader
	io.Writer
}

func (c fakeNetConn) Close() error                       { return nil }
func (c fakeNetConn) LocalAddr() net.Addr                { return localAddr }
func (c fakeNetConn) RemoteAddr() net.Addr               { return remoteAddr }
func (c fakeNetConn) SetDeadline(t time.Time) error      { return nil }
func (c fakeNetConn) SetReadDeadline(t time.Time) error  { return nil }
func (c fakeNetConn) SetWriteDeadline(t time.Time) error { return nil }

type fakeAddr int

var (
	localAddr  = fakeAddr(1)
	remoteAddr = fakeAddr(2)
)

func (a fakeAddr) Network() string {
	return "net"
}

func (a fakeAddr) String() string {
	return "str"
}

type testHijakableResponseWriter struct {
	brw *bufio.ReadWriter
	http.ResponseWriter
}

func (resp *testHijakableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return fakeNetConn{strings.NewReader(""), &bytes.Buffer{}}, resp.brw, nil
}
