package webserver

type Config struct {
	Enable            bool   `json:"enable"`
	Address           string `json:"address"`
	BasicAuthUser     string `json:"basic_auth_user"`
	BasicAuthPassword string `json:"basic_auth_password"`
	StaticDir         string `json:"static_dir"`
	EnableCORS        bool   `json:"enable_cors"`
}

var Cfg Config
