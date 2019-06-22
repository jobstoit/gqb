// Copyright 2019 Job Stoit. All rights reserved.

// gqb is made for generating a query builder.
//
// The gqb command takes the path to the db.yml as argument
//
// The command takes the following flags as arguments
//   -migrate 		Specifies the output for a generated sql migration
//
//   -db                Directly inserts the configured structure into the database
//                      using the DB_DRIVER and DB_CONNECTION_STRING enviroment
//                      variables or the flags for this mode
//
//   -dvr               Set the driver for the db flag
//
//   -cs                Sets the connection string for the db flag
//
//   -model             Writes the configuration to NiseVoid/qb model(s) takes
//                      the output file(s) as argument
//
//   -pkg               Used by the model flag, sets the package name of the
//                      model file(s)
//
package main

import "flag"

import "io/ioutil"
import "os"

// Flags for the application
// nolint: deadcode, varcheck, unused
var (
	migrateFlag = flag.String(`migrate`, ``, `Specifies the output for the SQL migration`)
	dbFlag      = flag.Bool(`db`, false, `Directly inserts the configured structur into the database`)
	dvrFlag     = flag.String(`dvr`, `postgres`, `Sets the driver for the db flag. (postgres by default)`)
	csFlag      = flag.String(`cs`, ``, `Sets the connection string for the db flag`)
	modelFlag   = flag.String(`model`, ``, `Sets the output for the Nisevoid/qb models`)
	pkgFlag     = flag.String(`pkg`, ``, `Used by the model flag, sets the package name of the model file(s)`)
)

var mdl Model

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	cFile := flag.Arg(0)
	if cFile == `` {
		flag.PrintDefaults()
		os.Exit(0)
	}

	bitz, err := ioutil.ReadFile(cFile)
	catch(err)
	mdl = readConfig(bitz)
}

func main() {
	if *modelFlag != `` {
		modelFile, err := os.Create(*modelFlag)
		catch(err)
		defer modelFile.Close()
		CreateQbModel(mdl, modelFile)
	}
}
