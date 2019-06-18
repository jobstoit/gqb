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

func quote(i ...interface{}) (s string) {
	s += "`"
	if len(i) == 0 {
		return
	}
	for e, it := range i {
		s += fmt.Sprint(it)
		if e != len(i)-1 {
			s += ` `
		}
	}
	s += "`"
	return
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
	{{title .Name}} qb.Field
{{- end}}
	table *qb.Table
}

// SQL is the qb.Query implementation for migration the {{title $t}} table
func (*{{title $t}}Type) SQL(_ qb.SQLBuilder) (q string, _ []interface{}) {
	q = {{quote}}CREATE TABLE IF NOT EXISTS {{print $t}} { {{range $n, $col := $cols}}
	{{$col.Name -}}
	{{- if notnil $col.DataType}} {{$col.DataType.Type}}{{end -}}
	{{- if gt $col.Size 0}}({{$col.Size}}){{end -}}
	{{- if $col.Primary}} PRIMARY{{end -}}
	{{- if $col.Unique}} UNIQUE{{end -}}
	{{- if not (eq $col.Default "")}} DEFAULT {{$col.Default}}{{end -}}
	{{- range $g, $c := $col.Constraints }} ADD CONSTRAINT {{$c}}{{end -}}
	{{- if $n}},{{end -}}
	{{- end -}}{{- quote}}
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

//insert starts a INSERT query
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
{{- else if gt (len $enu) 0}}
// {{title $t}}Type represents the enum "{{print $t}}
type {{title $t}}Type []string

// SQL is the qb.Query implementation for migrating the {{title $t}} enum
func (*{{title $t}}Type) SQL(_ qb.SQLBuilder) (q string, _[]interface{}) {
	q = {{quote}}CREATE ENUM {{$t}}
	{{- join $enu ", " -}}
	{{- quote}}
	return
}

// {{title $t}} returns a new {{title $t}}Type
func {{title $t}}() *{{title $t}}Type {
	return &{{title $t}}Type([]string{
		{{- join $enu ", " -}}
	})
}
{{end}}

{{end}}`
