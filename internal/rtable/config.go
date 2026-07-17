package rtable

type Config struct {
	Path          string `toml:"path"`
	DefaultAccept bool   `toml:"default_accept"`
}

func DefaultConfig() *Config {
	return &Config{
		Path:          "/etc/yfw/rule.json",
		DefaultAccept: true,
	}
}
