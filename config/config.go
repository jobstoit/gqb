// Copyright 2019 Job Stoit. All rights reserved.

package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// this config struct is for mapping the object to the nessasary state
type config struct {
	Driver string                       `yaml:"driver"`
	Pkg    string                       `yaml:"pkg"`
	Tables map[string]map[string]string `yaml:"tables,flow"`
	Enums  map[string][]string          `yaml:"enums,flow"`
}

// ReadConfig reads the data of the given yaml file into a model
func Read(data []byte) (m Model) {
	var c config
	catch(yaml.Unmarshal(data, &c), `Yaml configuration is unreadable`)

	defaultReg := regexp.MustCompile(`default\((\w+)\)`)
	constrainReg := regexp.MustCompile(`constraint\(([\w,\s;]+)\)`)

	m.Pkg = c.Pkg
	m.Driver = c.Driver
	m.Types = make(map[string]DataType)

	for i, tab := range c.Tables {
		for e, context := range tab {
			c := Column{}
			c.table = i
			c.Name = e

			if match := defaultReg.FindStringSubmatch(context); len(match) == 2 {
				c.Default = match[1]
			}

			if match := constrainReg.FindStringSubmatch(context); len(match) == 2 {
				c.Constraints = strings.Split(match[1], `;`)
			}

			typ, size := getRawType(context)
			c.rawType = typ
			c.Size = size
			c.Primary = strings.Contains(context, `primary`)
			c.Nullable = strings.Contains(context, `nullable`)
			c.Unique = strings.Contains(context, `unique`)

			key := i + `.` + e
			m.Types[key] = &c
		}
	}

	for i, enu := range c.Enums {
		m.Types[i] = &Enum{i, enu}
	}

	// Set the types of each column
	for _, typ := range m.Types {
		if col, ok := typ.(*Column); ok {
			m.GetType(col)
		}
	}

	return
}

var typeDataReg = regexp.MustCompile(`^([\w\_\.]+)(\(\s?(\d{0,3})\s?\))?`)

func getRawType(context string) (rawType string, size int) {
	typeData := typeDataReg.FindStringSubmatch(context)
	if len(typeData) == 0 {
		fatal(`Type not defined: ` + context)
	}

	rawType = typeData[1]
	if len(typeData) >= 4 && typeData[2] != `` {
		isize, err := strconv.Atoi(typeData[3])
		catch(err, `Datatype in an invalid format: %s\n`, rawType)
		size = isize
	}
	return
}

// Catch is used to panic a function/statement when errors occur
func catch(err error, msg string, args ...interface{}) {
	if err != nil {
		fatal(msg, args...)
	}
}

// Fatal closes the program with a message
func fatal(msg string, args ...interface{}) {
	if !strings.HasPrefix(msg, "\n") {
		msg += "\n"
	}

	fmt.Printf(msg, args...)
	os.Exit(1)
}
