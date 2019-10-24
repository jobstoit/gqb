// Copyright 2019 Job Stoit. All rights reserved.

package template

import (
	"fmt"
	"strings"
)

func title(i interface{}) (s string) {
	for _, part := range strings.Split(fmt.Sprint(i), `_`) {
		part = strings.Title(part)
		if part == `Id` ||
			part == `Url` ||
			part == `Sql` {
			part = strings.ToUpper(part)
		}
		s += part
	}
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

func notNil(i interface{}) bool {
	return i != nil
}
