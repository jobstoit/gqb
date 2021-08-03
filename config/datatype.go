// Copyright 2019 Job Stoit. All rights reserved.

package config

import "fmt"

// DataType is a models data structure
type DataType interface {
	Table() string
	Type() string
}

// Model is the full db model structure and configuration
// The model contains all the nessasary information for
// Creating Query builders
type Model struct {
	Driver string
	Pkg    string
	Types  map[string]DataType
}

// Tables returns all the tables in the model
func (x Model) Tables() (tables []string) {
	keys := make(map[string]bool)
	for _, t := range x.Types {
		if t.Table() != `` && !keys[t.Table()] {
			keys[t.Table()] = true
			tables = append(tables, t.Table())
		}
	}
	return
}

//TODO: think about transforming the Types over to a map in order to map them better
// also look if the Tables array is neccisary for what you want to have

// Column contains the table column properties and
// lies at the core of the query builder.
// Note that if the datatype is of type column that
// the column is a foreign key.
type Column struct {
	table       string
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

// Table is the DataType implementation
func (x Column) Table() string {
	return x.table
}

// Type is the DataType implementation
func (x Column) Type() string {
	if x.DataType == nil {
		fmt.Println(x.Name)
	}
	return x.DataType.Type()
}

func (x Column) String() string {
	return x.Table() + `.` + x.Name
}

// Enum is en e nummeric object which can be
// defined as type in the database
type Enum struct {
	table  string
	Values []string
}

// Table is the DataType implementation
func (x Enum) Table() string {
	return x.table
}

// Type is the DataType implementation
func (x Enum) Type() string {
	return getPrimitiveType(`varchar`).Type()
}

// PrimitiveType is a primative database type
type PrimitiveType string

// Table is the DataType implementation
func (x PrimitiveType) Table() string {
	return ``
}

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
			(c.rawType == col.Table()+`.`+col.Name ||
				c.rawType == col.Table()+`(`+col.Name+`)`) {
			c.DataType = col
			return
		} else if enum, ok := typ.(*Enum); ok && enum != nil &&
			c.rawType == enum.Table() {
			c.DataType = enum
			return
		}
	}
	fatal(`Type not found: %s\n`, c.rawType)
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

	case `time`, `datetime`:
		return getPrimitiveType(`timestamp`)

	case `bool`:
		return getPrimitiveType(`boolean`)

	case `varchar`, `text`, `int`, `tinyint`, `smallint`,
		`bigint`, `double`, `float`, `date`, `timestamp`, `boolean`:
		return PrimitiveType(i)
	default:
		return nil
	}
}
