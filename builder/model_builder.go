package builder

import (
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/goship/config"
	"github.com/gohouse/schema"
)

type Column struct {
	ColumnName    string
	ColumnType    string
	Tag           string
	ColumnComment string
}
type Method struct {
	MethodName    string
	ReturnType    string
	MethodContent string
}
type Table struct {
	TableName    string
	TableComment string
	TableColumns []Column
	TableMethods []Method
}

type Model struct {
	sm   *schema.Schema
	ge   *gorose.Engin
	conf *config.Config
}

func NewModel(ge *gorose.Engin, c *config.Config) *Model {
	return &Model{sm: schema.NewSchema(ge), ge: ge, conf: c}
}

func (m *Model) Build() {
	
}
