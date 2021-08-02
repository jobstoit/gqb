// Copyright 2019 Job Stoit. All rights reserved.

package model

import (
	"bytes"
	"strings"
	"testing"

	_ "embed"

	"github.com/jobstoit/gqb/config"
)

//go:embed db.model.yml
var testConfig []byte

func TestCreateQbModel(t *testing.T) {
	mdl := config.Read(testConfig)
	buff := new(bytes.Buffer)

	CreateQbModel(mdl, buff)
	res := buff.String()

	// Table
	if !strings.Contains(res, "qbUsersTable = qb.Table{Name: `users`}") {
		t.Error(`Error or wrong format table instantiation`)
	}

	// Columns
	if !strings.Contains(res, "qbUsersFID = qb.TableField{Parent: &qbUsersTable, Name: `id`, Type: qb.Int}") {
		t.Errorf("Error or wrong in table column format, maybe wrong type\n\n%s", res)
	}

	if !strings.Contains(res, "qbUsersFLastName = qb.TableField{Parent: &qbUsersTable, Name: `last_name`, Type: qb.String, Size: 100}") {
		t.Errorf("Error or wrong in table column format, maybe wrong or no Size\n\n%s", res)
	}

	if !strings.Contains(res, "qbUsersFBio = qb.TableField{Parent: &qbUsersTable, Name: `bio`, Type: qb.String, Nullable: true}") {
		t.Errorf("Error or wrong in table column format, maybe wrong nullable\n\n%s", res)
	}

	// enums
	if !strings.Contains(res, "type RolesType []string") {
		t.Errorf("Error or wrong in enum generation\n\n%s", res)
	}

}

func TestCreateMigration(t *testing.T) {
	mdl := config.Read(testConfig)
	buff := new(bytes.Buffer)

	CreateMigration(mdl, buff)
	res := buff.String()

	if !(strings.Contains(res, `CREATE TABLE IF NOT EXISTS users (`) &&
		strings.Contains(res, `id int NOT NULL PRIMARY KEY`) &&
		strings.Contains(res, `name varchar(100) NOT NULL`) &&
		strings.Contains(res, `last_name varchar(100) NOT NULL`) &&
		strings.Contains(res, `bio text`) &&
		strings.Contains(res, `role varchar DEFAULT 'GENERAL'`) &&
		strings.Contains(res, `);`)) {
		t.Errorf("Error or wrong in table query generation:\n\n%s", res)
	}

	if !strings.Contains(res, `CREATE TYPE roles AS ENUM ( ADMIN, GENERAL );`) {
		t.Errorf("Error or wrong in enum query generation:\n\n%s", res)
	}
}

func TestTitle(t *testing.T) {
	expectStr(t, `UserID`, title(`user_id`))
	expectStr(t, `UserModel`, title(`user_model`))
	expectStr(t, `ModelSQL`, title(`model_sql`))
	expectStr(t, `WebsiteURL`, title(`website_url`))
}

func TestQuote(t *testing.T) {
	expectStr(t, "`test`", quote(`test`))
	expectStr(t, "`furter testing`", quote(`furter testing`))
}

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
