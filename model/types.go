// Copyright 2019 Job Stoit. All rights reserved.

package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// generalNativeTypes
var generalNativeTypes = []DataType{
	NewNativeType(`varchar`, `string`, `char`, `character`, `charactering varying`),
	NewNativeType(`int`, `integer`),
	NewNativeType(`float`, `real`),
	NewNativeType(`timestamp`, `time`, `datetime`),
	NewNativeType(`boolean`, `bool`),
	NewNativeType(`text`),
	NewNativeType(`double`),
	NewNativeType(`date`),
	NewNativeType(`tinyint`),
	NewNativeType(`smallint`),
	NewNativeType(`bigint`),
}

// DataType is a type within the database context
// This can be either a native database type or
// a user specified type like "user.id"
type DataType interface {
	Type() string
	Refs() []string
}

// Table defines a table in the database
type Table string

// String is the fmt.Stringer implementation
func (t Table) String() string {
	return string(t)
}

// Column contains all the table collumn properties
// Note that the datatype of the collumn can be another
// collumn and if so is a foreign key otherwise it should
// rever to a native type
type Column struct {
	Table       *Table
	Name        string
	DataType    DataType
	Size        int
	Default     string
	Nullable    bool
	Unique      bool
	Primary     bool
	Constraints []string
	rawType     string
}

// NewColumn creates a new column given the configuration
func NewColumn(table *Table, name, context string) *Column {
	x := new(Column)
	x.Table = table
	x.Name = name

	// get type and size
	typeData := regexp.MustCompile(`^[a-zA-Z\.\_]+(\((\d{0,3})\))?`).
		FindStringSubmatch(context)
	if len(typeData) == 0 {
		panic(`syntax error in context. type not properly defined: ` + context)
	}

	x.rawType = typeData[0]
	if len(typeData) >= 2 {
		x.rawType = strings.Trim(x.rawType, typeData[1])
	}

	if len(typeData) == 3 && typeData[2] != `` {
		size, err := strconv.Atoi(typeData[2])
		if err != nil {
			panic(`invallid given type size`)
		}
		x.Size = size
	}

	// get default
	if match := regexp.MustCompile(`,\s?default\((\w+)\)`).
		FindStringSubmatch(context); len(match) == 2 {
		x.Default = match[1]
	}

	// get constraints
	if match := regexp.MustCompile(`,\s?constraint\(([\w,\s;]+)\)`).
		FindStringSubmatch(context); len(match) == 2 {
		x.Constraints = strings.Split(match[1], `;`)
	}

	// Get boolean types
	context = regexp.MustCompile(`\s`).ReplaceAllString(context, ``)
	x.Primary = strings.Contains(context, `,primary`)
	x.Nullable = strings.Contains(context, `,nullable`)
	x.Unique = strings.Contains(context, `,unique`)

	return x
}

// Type is a DataType implementation
func (x Column) Type() string {
	return x.DataType.Type()
}

// Refs is a DataType implementation
func (x Column) Refs() []string {
	return []string{
		fmt.Sprintf("%s.%s", x.Table, x.Name),
		fmt.Sprintf("%s(%s)", x.Table, x.Name),
	}
}

// NativeType is a native database type (like varchar, int, etc.)
type NativeType []string

// NewNativeType returns a native type which can create a
// Context for the database
func NewNativeType(name string, refs ...string) *NativeType {
	n := NativeType{name}
	n = append(n, refs...)
	return &n
}

// Type is the DataType implementation
func (x NativeType) Type() string {
	return []string(x)[0]
}

// Refs is a DataType implementation
func (x NativeType) Refs() []string {
	return []string(x)
}
