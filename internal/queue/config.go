package queue

type Config struct {
	Num  uint16 `toml:"num"`
	Name string `toml:"name"`
}

func DefaultConfig() *Config {
	return &Config{
		Num:  0,
		Name: "yfw",
	}
}
