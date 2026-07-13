package webserver

type Config struct {
	Enable       bool   `json:"enable"`
	Address      string `json:"address"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	EnableCORS   bool   `json:"enable_cors"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:       true,
		Address:      ":8080",
		AuthUsername: "admin",
		AuthPassword: "admin",
		EnableCORS:   true,
	}
}
