// Copyright 2019 Job Stoit. All rights reserved.

package template

import (
	"io"
	"strings"
	templ "text/template"

	"github.com/jobstoit/gqb/model"
)

// CreateQBModel creates the NiseVoid qb/model
func CreateQBModel(m model.Context, wr io.Writer) {
	err := templ.Must(templ.New(`NiseVoid/qb model`).
		Funcs(templ.FuncMap{
			`title`:  title,
			`quote`:  quote,
			`join`:   strings.Join,
			`notnil`: notNil,
			`qbType`: qbType,
		}).parse(queryTempl+qbTempl)).
		Execute(wr, m)

	if err != nil {
		panic(err)
	}
}

func qbType(x model.DataType) string {
	switch x.Type() {
	case `varchar`, `text`:
		return `qb.String`

	case `int`, `tinyint`, `smallint`, `bigint`:
		return `qb.Int`

	case `double`, `float`:
		return `qb.Float`

	case `date`, `datetime`:
		return `qb.Date`

	case `boolean`:
		return `qb.Bool`

	default:
		return `qb.Int`
	}
}

var qbTempl = `
{{range $e, $t := $.Tables -}}
{{- $cols := $t.Columns $.Types -}}
{{- $enu := $t.Enum $.Types -}}
// {{print $t}}
{{- if gt (len $cols) 0}}
var (
	qb{{title $t}}Table = qb.Table{Name: {{quote $t}}}
	{{range $cols}}
	qb{{title $t}}F{{title .Name}} = qb.TableField{Parent: &qb{{title $t}}Table, Name: {{quote .Name}}
	{{- if notnil .DataType}}, Type: {{qbtype .DataType}}{{end -}}
	{{- if gt .Size 0}}, Size: {{.Size}}{{end -}}
	{{- if .Nullable}}, Nullable: true{{end -}}
}{{end}}
)

// {{title $t}}Type represents the table "{{print $t}}"
type {{title $t}}Type struct {
{{- range $cols}}
	{{title .Name}} qb.Field{{end}}

	table *qb.Table
}

// SQL is the qb.Query implementation for migration the {{title $t}} table
func (*{{title $t}}Type) SQL(_ qb.SQLBuilder) (q string, _ []interface{}) {
	q = {{quote}}{{template "tablequery" $cols}}{{- quote}}
	return
}

// GetTable returns an object with info about the table
func (t *{{title $t}}Type) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *{{title $t}}Type) Select(f ...qb.Field) *qb.SelectBuilder {
	return t.table.Select(f)
}

// Delete creates a DELETE query
func (t *{{title $t}}Type) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts a UPDATE query
func (t *{{title $t}}Type) Update() *qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts a INSERT query
func (t *{{title $t}}Type) Insert(f ...qb.Field) *qb.InsertBuilder {
	return t.table.Insert(f)
}

// {{title $t}} returns a new {{title $t}}Type
func {{title $t}}() *{{title $t}}Type {
	table := qb{{title $t}}Table
	return &{{title $t}}Type{
	{{- range $cols}}
		qb{{title $t}}F{{title .Name}}.Copy(&table),
	{{- end}}
		&table,
	}
}
{{- else if notnil $enu.Table}}

// {{title $t}}Type represents the enum "{{print $t}}
type {{title $t}}Type []string

// SQL is the qb.Query implementation for migrating the {{title $t}} enum
func (*{{title $t}}Type) SQL(_ qb.SQLBuilder) (q string, _[]interface{}) {
	q = {{quote}}{{template "enumquery" $enu}}{{quote}}
	return
}

// {{title $t}} returns a new {{title $t}}Type
func {{title $t}}() *{{title $t}}Type {
	enu := {{title $t}}Type([]string{ {{range $enu.Values}}
		{{quote .}},{{end}}
	})
	return &enu
}
{{end}}

{{end}}`
