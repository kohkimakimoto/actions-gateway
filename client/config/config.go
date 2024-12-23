package config

import (
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the client configuration.
type Config struct {
	// This is the path to the config file.
	Path string `toml:"-"`
	// The server url that the client connects to.
	Server string `toml:"server"`
	// This is used as the server URL in the OpenAPI specification.
	// The default value is the same as the Server.
	// See also: https://swagger.io/docs/specification/v3_0/api-host-and-base-path/
	ServerApiURL string `toml:"server_api_url"`
	// This is the token for authentication.
	Token string `toml:"token"`
	// This is the directory that contains the actions.
	// If it is a relative path, it will be relative to the directory where the config file is located.
	ActionsDir string `toml:"actions_dir"`
	// This is the absolute path to the ActionsDir directory.
	ActionsAbsDir string `toml:"-"`
	// This is the status file path.
	// The status file is used to store the status of the client connection to the server.
	// It is also required when you connect to the server as a daemon.
	// If it is a relative path, it will be relative to the directory where the config file is located.
	StatusFile string `toml:"status_file"`
	// This is the absolute path to the StatusFile file.
	StatusAbsFile string `toml:"-"`
	// This is the pid file path.
	// The pid file is only used (and required) when you connect to the server as a daemon.
	// If it is a relative path, it will be relative to the directory where the config file is located.
	PidFile string `toml:"pid_file"`
	// This is the absolute path to the PidFile file.
	PidAbsFile string `toml:"-"`
	// This is the log file path.
	// The log file is only used (and required) when you connect to the server as a daemon.
	// If it is a relative path, it will be relative to the directory where the config file is located.
	LogFile string `toml:"log_file"`
	// This is the absolute path to the LogFile file.
	LogAbsFile string `toml:"-"`
	// This is the maximum number of reconnection attempts.
	// The default value is 10.
	MaxReconnectAttempts int `toml:"max_reconnect_attempts"`
	// This is the maximum backoff time in seconds.
	// The default value is 32.
	MaxReconnectBackoff int `toml:"max_reconnect_backoff"`
	// This is the info object config of the OpenAPI spec.
	// https://swagger.io/specification/#info-object
	SpecInfo *SpecInfoConfig `toml:"spec_info"`
}

// Dir returns the directory of the config file.
func (c *Config) Dir() string {
	return filepath.Dir(c.Path)
}

type SpecInfoConfig struct {
	Title       string `toml:"title"`
	Summary     string `toml:"summary"`
	Description string `toml:"description"`
	Version     string `toml:"version"`
}

func LoadFromFile(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	c := &Config{
		Path:     absPath,
		SpecInfo: &SpecInfoConfig{},
	}
	if _, err := toml.DecodeFile(c.Path, c); err != nil {
		return nil, err
	}

	if c.ServerApiURL == "" {
		c.ServerApiURL = c.Server
	}

	if c.ActionsDir == "" {
		c.ActionsDir = "actions"
	}

	var actionsAbsDir string
	if filepath.IsAbs(c.ActionsDir) {
		actionsAbsDir = c.ActionsDir
	} else {
		actionsAbsDir, err = filepath.Abs(filepath.Join(c.Dir(), c.ActionsDir))
		if err != nil {
			return nil, err
		}
	}
	c.ActionsAbsDir = actionsAbsDir

	if c.StatusFile != "" {
		var statusAbsFile string
		if filepath.IsAbs(c.StatusFile) {
			statusAbsFile = c.StatusFile
		} else {
			statusAbsFile, err = filepath.Abs(filepath.Join(c.Dir(), c.StatusFile))
			if err != nil {
				return nil, err
			}
		}
		c.StatusAbsFile = statusAbsFile
	}

	if c.PidFile != "" {
		var pidAbsFile string
		if filepath.IsAbs(c.PidFile) {
			pidAbsFile = c.PidFile
		} else {
			pidAbsFile, err = filepath.Abs(filepath.Join(c.Dir(), c.PidFile))
			if err != nil {
				return nil, err
			}
		}
		c.PidAbsFile = pidAbsFile
	}

	if c.LogFile != "" {
		var logAbsFile string
		if filepath.IsAbs(c.LogFile) {
			logAbsFile = c.LogFile
		} else {
			logAbsFile, err = filepath.Abs(filepath.Join(c.Dir(), c.LogFile))
			if err != nil {
				return nil, err
			}
		}
		c.LogAbsFile = logAbsFile
	}

	if c.MaxReconnectAttempts == 0 {
		c.MaxReconnectAttempts = 10
	}
	if c.MaxReconnectBackoff == 0 {
		c.MaxReconnectBackoff = 32
	}

	if c.SpecInfo.Title == "" {
		c.SpecInfo.Title = "Actions Gateway API"
	}
	if c.SpecInfo.Summary == "" {
		c.SpecInfo.Summary = "Actions Gateway API"
	}
	if c.SpecInfo.Description == "" {
		c.SpecInfo.Description = "Actions Gateway API"
	}
	if c.SpecInfo.Version == "" {
		c.SpecInfo.Version = "1.0.0"
	}

	return c, nil
}

var InitialConfig = strings.TrimLeft(`
# ------------------------------------------------------------
# This is a client config for Actions Gateway.
# See https://github.com/kohkimakimoto/actions-gateway
# ------------------------------------------------------------

# The server url that the client connects to.
# The following server is a free public server that is hosted by the author.
server = "https://actions-gateway.kohkimakimoto.dev"
# If you want to connect to your local server (for development purposes), you can use the following.
#server = "http://localhost:18800"

# This is used as the server URL in the OpenAPI specification.
# The default value is the same as the "server" config.
# See also: https://swagger.io/docs/specification/v3_0/api-host-and-base-path/
#server_api_url = "http://localhost:18800"

# This is a token for authentication.
# To get the token, you can use the 'actions-gateway new-token' command.
token = "replace-me-with-your-token..."

# This is the directory that contains the actions.
# If it is a relative path, it will be relative to the directory where the config file is located.
# The default value is "actions".
actions_dir = "actions"

# This is the status file path.
# The status file is used to store the status of the client connection to the server.
# It is also required when you connect to the server as a daemon.
# If it is a relative path, it will be relative to the directory where the config file is located.
status_file = "status.json"

# This is the pid file path.
# The pid file is only used (and required) when you connect to the server as a daemon.
# If it is a relative path, it will be relative to the directory where the config file is located.
pid_file = "client.pid"

# This is the log file path.
# The log file is only used (and required) when you connect to the server as a daemon.
# If it is a relative path, it will be relative to the directory where the config file is located.
log_file = "client.log"

# This is the maximum number of reconnection attempts.
# The default value is 10.
max_reconnect_attempts = 10

# This is the maximum backoff time in seconds.
# The default value is 32.
max_reconnect_backoff = 32

# ------------------------------------------------------------
# Spec info config.
# ------------------------------------------------------------

# This is the info object config of the OpenAPI spec.
# See more detail: https://swagger.io/specification/#info-object
# Currently, only the following fields are supported.
#spec_info.title = "Actions Gateway API"
#spec_info.summary = "Actions Gateway API"
#spec_info.description = "Actions Gateway API"
#spec_info.version = "1.0.0"

`, "\n")
