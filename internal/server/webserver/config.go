package webserver

type Config struct {
	Enable       bool   `toml:"enable"`
	Address      string `toml:"address"`
	AuthUsername string `toml:"auth_username"`
	AuthPassword string `toml:"auth_password"`
	StaticDir    string `toml:"static_dir"`
	EnableCORS   bool   `toml:"enable_cors"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:       true,
		// 只允许本地
		Address:      "0.0.0.0:8080",
		AuthUsername: "",
		AuthPassword: "",
		StaticDir:    "",
		EnableCORS:   true,
	}
}
