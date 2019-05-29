// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"io/ioutil"
	"testing"
)

func TestReadConfig(t *testing.T) {
	data, err := ioutil.ReadFile(`./db.test.yml`)
	catch(err)

	m := readConfig(data)
	expectInt(t, 5, len(m.Tables))

	for _, typ := range m.Types {
		if col, ok := typ.(*Column); ok {
			switch col.Name {
			case `role`: // test enum ref
				if fk, ok := col.DataType.(*Enum); !ok || fk.Values[0] != `ADMIN` {
					t.Error(`user.role does not refer to enum`)
				}
			case `created_by`: // test foreign key ref
				if fk, ok := col.DataType.(*Column); !ok ||
					fk.Table.String() != `user` ||
					fk.Name != `id` {
					t.Error(`user_id does not have the proper foreign key`)
				}
			case `created_at`: // test primitive type ref
				if col.DataType.Type() != `datetime` {
					t.Error(`created_at doesnt refer to datetime`)
				}
			case `subtitle`:
				if col.Nullable != true {
					t.Error(`subtitle shoul be nullable but isn't`)
				}
			case `title`: // test if default is set
				if col.Default != `new_post` {
					t.Errorf("default not properly set, expected '%s' but got '%s'\n", `new_post`, col.Default)
				}
			case `post_id`:
				if len(col.Constraints) != 2 {
					t.Errorf("post_id should have 2 constraints but has %d\n", len(col.Constraints))
				}
				if col.Constraints[0] != `index` {
					t.Error(`first post_id constraint should be 'index'`)
				}
			}
		}
	}
}
