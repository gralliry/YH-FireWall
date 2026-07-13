package webserver

import (
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/itable"
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

	GetConnections() []ctable.Info
	CloseConnection(id string) error

	GetInterfaces() []itable.Info
	GetProtocols() []string
}
