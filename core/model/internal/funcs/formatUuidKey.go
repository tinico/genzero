package funcs

import (
	"bytes"
	"fmt"

	"github.com/licat233/genzero/sql"
	"github.com/licat233/genzero/tools"
)

type FormatUuidKey struct {
	modelName string
	name      string
	req       string
	resp      string
	fullName  string
	Table     *sql.Table
}

var _ ModelFunc = (*FormatUuidKey)(nil)

func NewFormatUuidKey(t *sql.Table) *FormatUuidKey {
	modelName := modelName(t.Name)
	name := "formatUuidKey"
	req := "uuid string"
	resp := "string"
	fullName := fmt.Sprintf("%s(%s) %s", name, req, resp)
	return &FormatUuidKey{
		modelName: modelName,
		name:      name,
		resp:      resp,
		fullName:  fullName,
		Table:     t,
	}
}

func (t *FormatUuidKey) String() string {
	var buf = new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\nfunc (m *%s) %s {", t.modelName, t.fullName))
	buf.WriteString(fmt.Sprintf("\n\treturn fmt.Sprintf(\"cache:%s:uuid:%%v\", uuid)", tools.ToLowerCamel(t.Table.Name)))
	buf.WriteString("\n}\n")
	return buf.String()
}

func (s *FormatUuidKey) FullName() string {
	return s.fullName
}

func (s *FormatUuidKey) Req() string {
	return s.req
}

func (s *FormatUuidKey) Resp() string {
	return s.resp
}

func (s *FormatUuidKey) Name() string {
	return s.name
}

func (s *FormatUuidKey) ModelName() string {
	return s.modelName
}
