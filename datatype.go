// Copyright 2019 Job Stoit. All rights reserved.

package main

// DataType is a models data structure
type DataType interface {
	Type() string
}

// Model is the full db model structure and configuration
// The model contains all the nessasary information for
// Creating Query builders
type Model struct {
	Driver string
	Pkg    string
	Tables []*Table
	Types  []DataType
}

// Table is a definition of a database table
type Table string

// String is the stringer implementation
func (x Table) String() string {
	return string(x)
}

// Columns returns the columns with a reference to this specific table
func (x *Table) Columns(types []DataType) (c []*Column) {
	for _, typ := range types {
		if col, ok := typ.(*Column); ok && col.Table == x {
			c = append(c, col)
		}
	}
	return
}

// Enum returns the enum values if the type is an enum
func (x *Table) Enum(types []DataType) []string {
	for _, typ := range types {
		if enu, ok := typ.(*Enum); ok && enu.Table == x {
			return enu.Values
		}
	}
	return []string{}
}

// Query returns the table generate query
func (x *Table) Query(types []DataType) (q string) {
	col := x.Columns(types)
	enu := x.Enum(types)

	if len(col) > 0 {
		q += "CREATE TABLE " + x.String() + "{\n"
		q += ""
	}

	if len(enu) > 0 {
		q += "CREATE ENUM " + x.String() + "{\n"
	}

	return
}

// Column contains the table column properties and
// lies at the core of the query builder.
// Note that if the datatype is of type column that
// the column is a foreign key.
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

// Type is the datatype implementation
func (x Column) Type() string {
	return x.DataType.Type()
}

func (x Column) String() string {
	return string(*x.Table) + `.` + x.Name
}

// Enum is en e nummeric object which can be
// defined as type in the database
type Enum struct {
	Table  *Table
	Values []string
}

// Type is the DataType implementation
func (x Enum) Type() string {
	return getPrimitiveType(`string`).Type()
}

// PrimitiveType is a primative database type
type PrimitiveType string

// Type is the datatype implementation
func (x PrimitiveType) Type() string {
	return string(x)
}

// GetType returns the full type of the column
func (x Model) GetType(c *Column) {
	if p := getPrimitiveType(c.rawType); p != nil {
		c.DataType = p
		return
	}

	for _, typ := range x.Types {
		if col, ok := typ.(*Column); ok && col != nil &&
			(c.rawType == col.Table.String()+`.`+col.Name ||
				c.rawType == col.Table.String()+`(`+col.Name+`)`) {
			c.DataType = col
			return
		} else if enum, ok := typ.(*Enum); ok && enum != nil &&
			c.rawType == string(*enum.Table) {
			c.DataType = enum
			return
		}
	}
	panic(`Type not found`)
}

// getPrimativeType returns the primative type matching the
// given query
func getPrimitiveType(i string) DataType {
	switch i {
	case `string`, `char`, `character`, `charactering varying`:
		return getPrimitiveType(`varchar`)

	case `integer`:
		return getPrimitiveType(`int`)

	case `real`:
		return getPrimitiveType(`float`)

	case `time`, `timestamp`:
		return getPrimitiveType(`datetime`)

	case `bool`:
		return getPrimitiveType(`boolean`)

	case `varchar`, `text`, `int`, `tinyint`, `smallint`,
		`bigint`, `double`, `float`, `date`, `datetime`, `boolean`:
		return PrimitiveType(i)
	default:
		return nil
	}
}
