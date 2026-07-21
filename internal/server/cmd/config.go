package cmdserver

type Config struct {
	SocketPath string `toml:"socket_path"`
}

func DefaultConfig() *Config {
	return &Config{
		SocketPath: "/tmp/yfw.sock",
	}
}
