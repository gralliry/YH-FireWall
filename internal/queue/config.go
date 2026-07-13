package queue

type Config struct {
	No   uint16 `json:"no"`
	Name string `json:"name"`
}

func DefaultConfig() *Config {
	return &Config{
		No:   0,
		Name: "yfw",
	}
}
