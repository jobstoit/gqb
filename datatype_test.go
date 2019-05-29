// Copyright 2019 Job Stoit. All rights reserved.

package main

import "testing"

func TestGetPrimativeType(t *testing.T) {
	expectStr(t, `int`, getPrimitiveType(`integer`).Type())
	expectStr(t, `varchar`, getPrimitiveType(`char`).Type())
	expectStr(t, `float`, getPrimitiveType(`real`).Type())
	expectStr(t, `boolean`, getPrimitiveType(`bool`).Type())
	expectStr(t, `datetime`, getPrimitiveType(`time`).Type())
	expectStr(t, `smallint`, getPrimitiveType(`smallint`).Type())
}

func TestGetType(t *testing.T) {
	for _, typ := range testModel.Types {
		if col, ok := typ.(*Column); ok {
			testModel.GetType(col)
			if col.Name == `user_id` && col.Type() != `int` {
				t.Error(`user_id is not of the expected type`)
			}
			if col.Name == `role` && col.Type() != `varchar` {
				t.Errorf("role does not refer to role enum: %s", col.Type())
			}
		}
	}

}

func TestTable(t *testing.T) {
	amm := 0
	for _, i := range testTableUser.Columns(testModel.Types) {
		switch i.String() {
		case `user.id`, `user.name`, `user.role`:
			amm++
		default:
			break
		}
	}
	expectInt(t, 3, amm)
	expectInt(t, 2, len(testTableRole.Enum(testModel.Types)))
	expectInt(t, 0, len(testTableUser.Enum(testModel.Types)))
}
