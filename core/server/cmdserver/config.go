package cmdserver

type Config struct {
	Enable     bool   `json:"enable"`
	SocketPath string `json:"socket_path"`
}

var Cfg Config
