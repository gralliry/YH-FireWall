package cmdserver

type Config struct {
	SocketPath string `json:"socket_path"`
}

func DefaultConfig() *Config {
	return &Config{
		SocketPath: "/tmp/yfw.sock",
	}
}
