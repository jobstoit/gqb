// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jobstoit/gqb/config"
)

func initMainTest() {
	if bitz, err := ioutil.ReadFile(`./db.test.yml`); err == nil {
		mdl = config.Read(bitz)
	} else {
		panic(`unable to read the db.test.yml file`)
	}

	dvr := `postgres`
	pkg := `main`
	model := `out/test.mdl.go`
	migration := `out/migration.sql`

	dvrFlag = &dvr
	pkgFlag = &pkg
	modelFlag = &model
	migrateFlag = &migration
}

func TestMain(m *testing.M) {
	initMainTest()

	os.Exit(m.Run())
}

func TestGeneral(t *testing.T) {
	initMainTest()

	main()
}
