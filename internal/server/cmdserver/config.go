package cmdserver

type Config struct {
	Enable     bool   `json:"enable"`
	SocketPath string `json:"socket_path"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:     true,
		SocketPath: "/tmp/yfw.sock",
	}
}
