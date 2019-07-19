// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
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
func ReadConfig(data []byte) (m Model) {
	c := config{}
	catch(yaml.Unmarshal(data, &c))

	removeSpacesReg := regexp.MustCompile(`\s`)
	defaultReg := regexp.MustCompile(`,\s?default\((\w+)\)`)
	constrainReg := regexp.MustCompile(`,\s?constraint\(([\w,\s;]+)\)`)

	m.Pkg = c.Pkg
	m.Driver = c.Driver

	for i, tab := range c.Tables {
		t := Table(i)
		for e, context := range tab {
			c := Column{}
			c.Table = &t
			c.Name = e

			if match := defaultReg.FindStringSubmatch(context); len(match) == 2 {
				c.Default = match[1]
			}

			if match := constrainReg.FindStringSubmatch(context); len(match) == 2 {
				c.Constraints = strings.Split(match[1], `;`)
			}

			context = removeSpacesReg.ReplaceAllString(context, ``)
			c.rawType, c.Size = getRawType(context)
			c.Primary = strings.Contains(context, `,primary`)
			c.Nullable = strings.Contains(context, `,nullable`)
			c.Unique = strings.Contains(context, `,unique`)

			m.Types = append(m.Types, &c)
		}
		m.Tables = append(m.Tables, &t)
	}

	for i, enu := range c.Enums {
		t := Table(i)
		m.Tables = append(m.Tables, &t)
		m.Types = append(m.Types, &Enum{&t, enu})
	}

	for _, typ := range m.Types {
		if col, ok := typ.(*Column); ok {
			m.GetType(col)
		}
	}

	return
}

var typeDataReg = regexp.MustCompile(`^[\w\_\.]+(\(\d{0,3}\))?`)

func getRawType(context string) (rawType string, size int) {
	typeData := typeDataReg.FindStringSubmatch(context)
	if len(typeData) == 0 {
		panic(`type not defined: ` + context)
	}

	rawType = typeData[0]
	if len(typeData) >= 2 && typeData[1] != `` {
		ssize := strings.Trim(typeData[1], `(`)
		ssize = strings.Trim(ssize, `)`)
		isize, err := strconv.Atoi(ssize)
		catch(err)

		size = isize
		rawType = strings.Trim(rawType, typeData[1])
	}
	return
}
