package config

type Queue struct {
	Num    uint16 `json:"num"`
	Accept bool   `json:"accept"`
}
