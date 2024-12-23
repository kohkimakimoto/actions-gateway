package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kohkimakimoto/actions-gateway/client/actions"
	"github.com/kohkimakimoto/actions-gateway/client/config"
	"github.com/kohkimakimoto/actions-gateway/client/status"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/kohkimakimoto/actions-gateway/version"
)

// Client is the client used to communicate with Actions Gateway server.
type Client struct {
	// config is the client configuration.
	config *config.Config
	// writer is the writer used to write messages.
	writer io.Writer
	// errWriter is the writer used to write error messages.
	errWriter io.Writer
	// httpClient is the HTTP client used to communicate with Actions Gateway server.
	httpClient *http.Client
	// reconnectAttempts is the current state of the reconnect attempts.
	reconnectAttempts int
}

// New creates a new client instance.
func New(cfg *config.Config, w io.Writer, errW io.Writer) *Client {
	return &Client{
		config:            cfg,
		writer:            w,
		errWriter:         errW,
		httpClient:        &http.Client{},
		reconnectAttempts: 0,
	}
}

func (c *Client) NewToken() (string, error) {
	resp, err := c.post("/api/new-token", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a new token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create a new token with server status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the response body: %w", err)
	}

	ret := &types.NewTokenResponse{}
	if err := json.Unmarshal(body, ret); err != nil {
		return "", fmt.Errorf("failed to parse the response body: %w", err)
	}

	return ret.Token, nil
}

func (c *Client) Connect(m *actions.ActionManager, sw *status.Writer) (err error) {
	maxBackoff := time.Duration(c.config.MaxReconnectBackoff) * time.Second
	for c.reconnectAttempts < c.config.MaxReconnectAttempts {
		if err = c.connect(m, sw); err != nil {
			c.reconnectAttempts++
			backoff := time.Duration(1<<c.reconnectAttempts) * time.Second
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			_, _ = fmt.Fprintf(c.errWriter, "Failed to connect: %v\n", err)
			_, _ = fmt.Fprintf(c.errWriter, "Reconnecting after %s (attempt %d/%d)\n", backoff, c.reconnectAttempts, c.config.MaxReconnectAttempts)

			time.Sleep(backoff)
			continue
		}
		break
	}

	// ignore the error for preventing overwriting the error
	_ = sw.UpdateToInactive(err)

	return err
}

func (c *Client) connect(m *actions.ActionManager, sw *status.Writer) error {
	spec, err := m.OutputSpec(c.errWriter)
	if err != nil {
		return fmt.Errorf("failed to make the spec: %w", err)
	}

	sessionNewRequest := &types.SessionNewRequest{
		Actions: m.ActionNames(),
		Spec:    spec,
	}

	if err := sw.UpdateToConnecting(sessionNewRequest); err != nil {
		return fmt.Errorf("failed to update the status: %w", err)
	}

	// request to the server to create a new session
	resp, err := c.post("/api/session/new", sessionNewRequest)
	if err != nil {
		return fmt.Errorf("failed to create a new session: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := &types.ErrorResponse{}
		if err := json.Unmarshal(body, errResponse); err != nil {
			return fmt.Errorf("failed to create a new session: statusCode=%s, body=%s", resp.Status, body)
		} else {
			return fmt.Errorf("failed to create a new session: statusCode=%s, error=%s", resp.Status, errResponse.Error)
		}
	}

	sessionNewResponse := &types.SessionNewResponse{}
	if err := json.Unmarshal(body, sessionNewResponse); err != nil {
		return fmt.Errorf("failed to parse the response body: %w", err)
	}

	// handle the interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(interrupt)
		close(interrupt)
	}()

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.Token))
	header.Add("User-Agent", fmt.Sprintf("actions-gateway/%s", version.Version))

	// connect to the websocket server
	conn, _, err := websocket.DefaultDialer.Dial(sessionNewResponse.URL, header)
	if err != nil {
		return fmt.Errorf("failed to connect to the websocket server: %w", err)
	}
	defer conn.Close()

	// reset the reconnect attempts, because the connection is successful
	c.reconnectAttempts = 0
	_, _ = fmt.Fprintf(c.writer, "Connected to the server: %s\n", c.makeURL("/"))

	// update status
	if err := sw.UpdateToActive(sessionNewRequest, sessionNewResponse); err != nil {
		return fmt.Errorf("failed to update the status: %w", err)
	}

	// channel to handle closing the connection
	done := make(chan struct{})

	// goroutine to receive messages from the server
	go func() {
		defer close(done)
		for {
			// read a message from the server
			_, message, err := conn.ReadMessage()
			if err != nil {
				// This error means disconnection from the server.
				return
			}
			go c.handleActionMessage(m, message)
		}
	}()

	// main loop
	for {
		select {
		case <-done:
			return errors.New("server disconnected")
		case sig := <-interrupt:
			// handle the interrupt signal
			_, _ = fmt.Fprintf(c.writer, "Received signal (%s).\n", sig.String())
			// send a close message to the server
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("failed to send a close message to the server: %w", err)
			}

			// wait for the server to close the connection or timeout
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func (c *Client) handleActionMessage(m *actions.ActionManager, message []byte) {
	_, _ = fmt.Fprintf(c.writer, "Received message: %s\n", message)

	// setup a result object
	result := &types.ActionResult{}

	// parse the message
	msg := &types.ActionMessage{}
	if err := json.Unmarshal(message, msg); err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "Failed to parse the message: %v\n", err)
		return
	}
	result.Id = msg.Id

	action := m.GetAction(msg.Name)
	if action == nil {
		_, _ = fmt.Fprintf(c.errWriter, "Failed to find the action: %s\n", msg.Name)
		result.Status = types.ActionResultStatusError
		result.Body = `{"error": "action not found"}`
		if err := c.NotifyResult(result); err != nil {
			_, _ = fmt.Fprintf(c.writer, "Failed to notify the result: %v\n", err)
		}
		return
	}
	b, err := actions.NewActionRunner(action, c.config.Dir(), c.errWriter).Run(msg)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "Failed to run the action: %v\n", err)
		result.Status = types.ActionResultStatusError
	} else {
		result.Status = types.ActionResultStatusSuccess
	}

	if b != nil {
		result.Body = string(b)
	}

	if err := c.NotifyResult(result); err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "Failed to notify the result: %v\n", err)
	}
}

func (c *Client) NotifyResult(result *types.ActionResult) error {
	resp, err := c.post("/api/notify", result)
	if err != nil {
		return fmt.Errorf("failed to notify the action result: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) post(urlPath string, payload any) (*http.Response, error) {
	var payloadBytes []byte
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to convert payload to JSON: %w", err)
		}
		payloadBytes = b
	}

	req, err := http.NewRequest(http.MethodPost, c.makeURL(urlPath), bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %w", err)
	}
	c.setHeaders(req)

	return c.httpClient.Do(req)
}

func (c *Client) makeURL(urlPath string) string {
	// concat the server URL and the URL path
	return fmt.Sprintf("%s%s", strings.TrimSuffix(c.config.Server, "/"), path.Clean(urlPath))
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", fmt.Sprintf("actions-gateway/%s", version.Version))
	if c.config.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.Token))
	}
}
