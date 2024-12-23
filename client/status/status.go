package status

import "github.com/kohkimakimoto/actions-gateway/server/types"

type CodeType string

const (
	CodeInactive   CodeType = "inactive"
	CodeConnecting CodeType = "connecting"
	CodeActive     CodeType = "active"
)

type Status struct {
	StatusCode         CodeType                  `json:"status_code"`
	SessionNewRequest  *types.SessionNewRequest  `json:"session_new_request,omitempty"`
	SessionNewResponse *types.SessionNewResponse `json:"session_new_response,omitempty"`
	Error              string                    `json:"error,omitempty"`
}

var initialStatus = &Status{
	StatusCode: CodeInactive,
}
