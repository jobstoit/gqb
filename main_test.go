// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func initMainTest() {
	if bitz, err := ioutil.ReadFile(`./db.test.yml`); err == nil {
		mdl = ReadConfig(bitz)
	} else {
		os.Exit(1)
	}

	flag.Set(`dvr`, `postgres`)
	flag.Set(`pkg`, `main`)
	flag.Set(`model`, `out/qb.mdl.go`)
	flag.Set(`migration`, `out/migration.go`)
	if !flag.Parsed() {
		flag.Parse()
	}
}

func TestMain(m *testing.M) {
	initMainTest()

	os.Exit(m.Run())
}

func TestGeneral(t *testing.T) {
	initMainTest()

	main()
}
