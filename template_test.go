// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCreateQbModel(t *testing.T) {
	buff := bytes.NewBufferString(``)
	CreateQbModel(testModel, buff)
	res := buff.String()

	// Table
	if !strings.Contains(res, "qbUserTable = qb.Table{Name: `user`}") {
		t.Fatal(`Error or wrong format table instantiation`)
	}

	// Columns
	if !strings.Contains(res, "qbUserFID = qb.TableField{Parent: &qbUserTable, Name: `id`, Type: qb.Int}") {
		t.Fatalf("Error or wrong in table column format, maybe wrong type\n\n%s", res)
	}

	if !strings.Contains(res, "qbUserFLastName = qb.TableField{Parent: &qbUserTable, Name: `last_name`, Type: qb.String, Size: 100}") {
		t.Fatalf("Error or wrong in table column format, maybe wrong or no Size\n\n%s", res)
	}

	if !strings.Contains(res, "qbUserFBio = qb.TableField{Parent: &qbUserTable, Name: `bio`, Type: qb.String, Nullable: true}") {
		t.Fatalf("Error or wrong in table column format, maybe wrong nullable\n\n%s", res)
	}

	// enums
	if !strings.Contains(res, "type RoleType []string") {
		t.Fatalf("Error or wrong in enum generation\n\n%s", res)
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
