// Copyright 2019 Job Stoit. All rights reserved.

package model

import (
	"bytes"
	"testing"
)

func TestReadAndCreateContext(t *testing.T) {
	c := `# This is a test configuration

tables:
  users:
    id: int, primary
    name: varchar(100)
    profile_photo: media.id
    bio: text, nullable
    role: string, default(GENERAL)
    ammount: real
    active: bool
    
  media:
    id: int, primary
    user_id: users.id
    path: varchar, unique
    order: integer, constraint(INDEX;NOT NULL)
    created_at: datetime, default(NOW)
`
	ctx := FromFile(bytes.NewBufferString(c))
	if l := len(ctx.Tables); l != 2 {
		t.Errorf(`Not enough tables in context want 2 have %d`, l)
	}

	for _, typ := range ctx.Types {
		if col, ok := typ.(*Column); ok {
			switch col.Name {
			case `id`:
				eq(t, true, col.Primary, `id not primary`)
				eq(t, false, col.Nullable, `id is nullable`)

			case `name`:
				eq(t, `varchar`, col.DataType.Type(), `name has not varchar as type`)
				eq(t, 100, col.Size, `name has unexpected size`)

			case `bio`:
				eq(t, `text`, col.DataType.Type(), `users.bio type is wrong`)
				eq(t, true, col.Nullable, `users.bio should be nullable`)

			case `role`:
				eq(t, `varchar`, col.DataType.Type(), `users.role should be varchar`)
				eq(t, `GENERAL`, col.Default, `users.role should have default`)

			case `user_id`:
				eq(t, `int`, col.DataType.Type(), `media.user_id should be of type int`)
				if fk, ok := col.DataType.(*Column); !ok ||
					fk.Refs()[0] != `users.id` {
					t.Errorf("foreign key does not refer to foreign object: %s", col.Name)
				}

			case `path`:
				eq(t, true, col.Unique, `media.path should have unique property`)

			case `order`:
				eq(t, 2, len(col.Constraints), `media.order should have contraints`)
				eq(t, `INDEX`, col.Constraints[0])
				eq(t, `NOT NULL`, col.Constraints[1])
			}
		}
	}

	c = `# This testfile tests if the verification fails with false types
tables:
  user:
    id: int, primary
    weirdtype: this_type_should_fail 
`
	defer func() {
		if r := recover(); r == nil {
			t.Error(`given configuration should panic but doesnt`)
		}
	}()
	FromFile(bytes.NewBufferString(c))
}

func eq(t *testing.T, expected, actual interface{}, message ...string) {
	if actual != expected {
		if len(message) != 0 {
			t.Errorf("TEST FAIL:\nexpected: %v\nactual: %v\n%s\n", expected, actual, message[0])
		} else {
			t.Errorf("TEST FAIL:\nexpected: %v\nactual: %v\n", expected, actual)
		}
	}
}
