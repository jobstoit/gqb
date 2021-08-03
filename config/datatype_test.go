// Copyright 2019 Job Stoit. All rights reserved.

package config

import "testing"

func TestGetPrimativeType(t *testing.T) {
	expectStr(t, `int`, getPrimitiveType(`integer`).Type())
	expectStr(t, `varchar`, getPrimitiveType(`char`).Type())
	expectStr(t, `float`, getPrimitiveType(`real`).Type())
	expectStr(t, `boolean`, getPrimitiveType(`bool`).Type())
	expectStr(t, `timestamp`, getPrimitiveType(`time`).Type())
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
