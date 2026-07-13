package rtable

type Config struct {
	Path          string `json:"path"`
	DefaultAccept bool   `json:"default_accept"`
}

func DefaultConfig() *Config {
	return &Config{
		Path:          "/etc/yfw/rule.json",
		DefaultAccept: true,
	}
}
