package rtable

type Config struct {
	Path          string `json:"path"`
	DefaultAccept bool   `json:"default_accept"`
}

var Cfg Config
