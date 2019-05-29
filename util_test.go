// Copyright 2019 Job Stoit. All rights reserved.

package main

import "testing"

var (
	testTableUser  = Table(`user`)
	testTableImage = Table(`image`)
	testTableRole  = Table(`role`)

	testModel = Model{
		Tables: []*Table{&testTableUser, &testTableImage, &testTableRole},
		Types: []DataType{
			&Column{
				Table:   &testTableUser,
				Name:    `id`,
				Primary: true,
				rawType: `int`,
			},
			&Column{
				Table:   &testTableUser,
				Name:    `name`,
				rawType: `varchar`,
				Size:    100,
			},
			&Column{
				Table:   &testTableUser,
				Name:    `last_name`,
				rawType: `varchar`,
				Size:    100,
			},
			&Column{
				Table:    &testTableUser,
				Name:     `bio`,
				rawType:  `text`,
				Nullable: true,
			},
			&Column{
				Table:   &testTableUser,
				Name:    `role`,
				rawType: `role`,
			},
			&Column{
				Table:   &testTableImage,
				Name:    `id`,
				Primary: true,
				rawType: `int`,
			},
			&Column{
				Table:   &testTableImage,
				Name:    `user_id`,
				rawType: `user.id`,
			},
			&Enum{
				Table:  &testTableRole,
				Values: []string{`ADMIN`, `GENERAL`},
			},
		},
	}
)

func expectInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("TEST FAILED\nExpected: %d\nActual: %d\n", expected, actual)
	}
}

func expectStr(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("TEST FAILED\nExpected: %s\nActual: %s\n", expected, actual)
	}
}
