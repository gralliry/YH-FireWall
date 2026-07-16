package queue

type Config struct {
	Num  uint16 `json:"no"`
	Name string `json:"name"`
}

func DefaultConfig() *Config {
	return &Config{
		Num:  0,
		Name: "yfw",
	}
}
