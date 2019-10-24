// Copyright 2019 Job Stoit. All rights reserved.

package template

import (
	"io"
	templ "text/template"

	"github.com/jobstoit/gqb/model"
)

// CreateInitialMigration generates an sql initial migration
func CreateInitialMigration(m model.Context, wr io.Writer) {
	err := templ.Must(templ.New(`initial-migration`).
		Funcs(templ.FuncMap{
			`notnil`: notNil,
			`column`: column,
		}).
		Parse(colTemplate, migrationTempl)).
		Execute(wr, m)

	if err != nil {
		panic(err)
	}
}

func column(t model.DataType) (c *model.Column) {
	if ty, ok := t.(*Column); ok {
		return ty
	}
	return new(model.Column)
}

var migrationTempl = `-- Tables
{{range $.Tables}}
CREATE TABLE IF NOT EXISTS {{print .}};
{{- end}}

-- Columns
{{range .Types}}
{{template "columnquery" (column .)}}
{{end}}`

var colTempl = `{{define "columnquery""}}ALTER TABLE {{print .Table}}
ADD {{ .Name }} {{.DataType.Type -}}
{{- if gt .Size 0}}({{.Size}}){{end -}
{{- if and (not .Nullable) (eq .Default "")}} NOT NULL{{end -}}
{{- if .Primary}} PRIMARY KEY{{end -}}
{{- if .Unique}} UNIQUE{{end -}}
{{- if not (eq .Default "")}} DEFAULT '{{.Default}}'{{end -}}
{{- range  $g, $c := .Constraints}}{{if $g}}, {{end}} ADD CONSTRAINT {{$c}}{{end}};{{end}}`
