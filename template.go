// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

// This file contains the templates for the code generation
// These are the template for each table

// CreateQbModel creates the models template
func CreateQbModel(m Model, wr io.Writer) {
	temp := template.Must(template.New(`qb-model`).
		Funcs(template.FuncMap{
			`title`:  title,
			`quote`:  quote,
			`join`:   strings.Join,
			`cols`:   cols,
			`enums`:  enums,
			`iscol`:  isCol,
			`isenum`: isEnum,
			`notnil`: notNil,
			`qbtype`: qbType,
		}).
		Parse(tableTempl))
	catch(temp.Execute(wr, m))
}

func title(s interface{}) (t string) {
	for _, part := range strings.Split(fmt.Sprint(s), `_`) {
		switch part {
		case `id`, `Id`, `sql`, `Sql`, `url`, `Url`:
			t += strings.ToUpper(part)
		default:
			t += strings.Title(part)
		}
	}
	return
}

func quote(i interface{}) string {
	return "`" + fmt.Sprint(i) + "`"
}

func cols(t *Table, types []DataType) []*Column {
	return t.Columns(types)
}

func isCol(t *Table, types []DataType) bool {
	return len(t.Columns(types)) > 0
}

func enums(t *Table, types []DataType) string {
	return strings.Join(t.Enum(types), `", "`)
}

func isEnum(t *Table, types []DataType) bool {
	return len(t.Enum(types)) > 0
}

func notNil(x interface{}) bool {
	return x != nil
}

func qbType(x DataType) string {
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

var tableTempl = `{{range $e, $t := $.Tables -}}
// {{print $t}}
{{- if iscol $t $.Types}}
var (
	qb{{title $t}}Table = qb.Table{Name: {{quote $t}}}
	{{range cols $t $.Types}}
	qb{{title $t}}F{{title .Name}} = qb.TableField{Parent: &qb{{title $t}}Table, Name: {{quote .Name}}
	{{- if notnil .DataType}}, Type: {{qbtype .DataType}}{{end -}}
	{{- if gt .Size 0}}, Size: {{.Size}}{{end -}}
	{{- if .Nullable}}, Nullable: true{{end -}}
}{{end}}
)

// {{title $t}}Type represents the table "{{print $t}}"
type {{title $t}}Type struct {
{{- range cols $t $.Types}}
	{{title .Name}} qb.Field
{{- end}}
	table *qb.Table
}

// SQL is the qb.Query implementation for migration the {{title $t}} table
func (*{{title $t}}Type) SQL(_ qb.SQLBuilder) (q string, _ []interface{}) {
	q = {{quote .Query}}
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
	{{- range cols $t $.Types}}
		qb{{title $t}}F{{title .Name}}.Copy(&table),
	{{- end}}
		&table,
	}
}
{{- else if isenum $t $.Types}}

// {{title $t}}Type represents the enum "{{print $t}}
type {{title $t}}Type []string

// SQL is the qb.Query implementation for migrating the {{title $t}} enum
func (*{{title $t}}Type) SQL(_ qb.SQLBuilder) (q string, _[]interface{}) {
	q = {{quote .Query}}
	return
}

// {{title $t}} returns a new {{title $t}}Type
func {{title $t}}() *{{title $t}}Type {
	enu := {{title $t}}Type([]string{
		"{{- enums $t $.Types -}}",
	})
	return &enu
}
{{end}}

{{end}}`
