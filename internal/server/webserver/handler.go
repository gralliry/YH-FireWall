package webserver

import (
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/rule"
)

type Handler interface {
	AppendRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	SearchRules() []rule.Info
	EnableRule(id string, enable bool) bool

	GetConfig() string
	SetConfig(raw string) error

	GetConnections() []conn.Info
	CloseConnection(id string) error

	GetInterfaces() []itf.Itf
	GetProtocols() []string
}
