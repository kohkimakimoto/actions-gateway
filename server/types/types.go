package types

type NewTokenResponse struct {
	// Token is a JWT token
	Token string `json:"token"`
}

type SessionNewRequest struct {
	// Actions is a list of action names
	Actions []string `json:"actions"`
	// Spec is a OpenAPI spec in YAML format that is supported by the session.
	Spec string `json:"spec"`
}

type SessionNewResponse struct {
	URL string `json:"url"`
}

type ActionMessage struct {
	// Id is a unique identifier for the action message
	Id string `json:"id"`
	// Name is the action name
	Name string `json:"name"`
	// Body is a payload of the action
	Body string `json:"body"`
}

type ActionResultStatus string

const (
	ActionResultStatusSuccess ActionResultStatus = "success"
	ActionResultStatusError   ActionResultStatus = "error"
)

type ActionResult struct {
	// It is the same as the id of the action message
	Id string `json:"id"`
	// "success" or "error"
	Status ActionResultStatus `json:"status"`
	// The result of the action that is produced from the action's STDOUT
	Body string `json:"error"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
