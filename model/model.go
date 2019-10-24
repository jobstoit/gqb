// Copyright 2019 Job Stoit. All rights reserved.

// Package model contains the go implementation for a database context
//
package model

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v2"
)

// Context consists of the database context with all
// the specified tables and datatypes
type Context struct {
	Tables []*Table
	Types  []DataType
}

// NewContext creates a new context
func NewContext(nativeTypes ...DataType) *Context {
	x := new(Context)

	if len(nativeTypes) == 0 {
		nativeTypes = generalNativeTypes
	}

	x.Types = append(x.Types, nativeTypes...)

	return x
}

// FromFile reads a yaml configuration and puts it in context
func FromFile(r io.Reader, native ...DataType) *Context {
	var conf struct {
		Tables map[string]map[string]string `yaml:"tables,flow"`
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		panic(`file not readable`)
	}

	if err := yaml.Unmarshal(buf.Bytes(), &conf); err != nil {
		panic(`syntax error in config, can not read yaml`)
	}

	ctx := NewContext(native...)
	ctx.AddTables(conf.Tables)

	return ctx
}

// AddTables adds tables and verifies the context
func (x *Context) AddTables(tables map[string]map[string]string) {

	// Filling the context
	for tableName, tableContent := range tables {
		t := Table(tableName)
		for columnName, columnContext := range tableContent {
			x.Types = append(x.Types, NewColumn(&t, columnName, columnContext))
		}
		x.Tables = append(x.Tables, &t)
	}

	// Setting the types and verifying
	for _, typ := range x.Types {
		if col, ok := typ.(*Column); ok {
			col.DataType = getType(x, col.rawType)
		}
	}
}

// getNativeType returns a native database type if given input resembels
// a natvie type
func getType(cxt *Context, i string) DataType {
	for _, typ := range cxt.Types {
		for _, ref := range typ.Refs() {
			if ref == i {
				return typ
			}
		}
	}
	panic(`type not found: ` + i)
}
