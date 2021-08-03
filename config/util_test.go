package config

import "testing"

var (
	testTableUser  = `user`
	testTableImage = `image`
	testTableRole  = `role`

	testModel = Model{
		Types: map[string]DataType{
			testTableUser + `.id`: &Column{
				table:   testTableUser,
				Name:    `id`,
				Primary: true,
				rawType: `int`,
			},
			testTableUser + `.name`: &Column{
				table:   testTableUser,
				Name:    `name`,
				rawType: `varchar`,
				Size:    100,
			},
			testTableUser + `.last_name`: &Column{
				table:   testTableUser,
				Name:    `last_name`,
				rawType: `varchar`,
				Size:    100,
			},
			testTableUser + `.bio`: &Column{
				table:    testTableUser,
				Name:     `bio`,
				rawType:  `text`,
				Nullable: true,
			},
			testTableUser + `.role`: &Column{
				table:   testTableUser,
				Name:    `role`,
				rawType: `role`,
				Default: `GENERAL`,
			},
			testTableImage + `.id`: &Column{
				table:   testTableImage,
				Name:    `id`,
				Primary: true,
				rawType: `int`,
			},
			testTableImage + `user_id`: &Column{
				table:   testTableImage,
				Name:    `user_id`,
				rawType: `user.id`,
			},
			testTableRole: &Enum{
				table:  testTableRole,
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
