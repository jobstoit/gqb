// Copyright 2019 Job Stoit. All rights reserved.

// gqb is made for generating a query builder.
//
// The gqb command takes the path to the db.yml as argument.
// So for example:
//   gqb -model ./models/qb.mdl.go db.yml
//
// The command takes the following flags as arguments.
//   -migrate           Specifies the output for a generated sql migration
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
// Configuration yaml example.
//   pkg: models                              # optional, default model
//   driver: postgres                         # optional, default postgres
//   tables:
//     users:
//       id: int, primary
//       email: varchar, unique
//       password: varhcar
//       fist_name: varchar(50)
//       last_name: varchar(100)
//       role: role                           # foreign key (enum)
//
//     posts:
//       id: int, primary
//       created_by: users.id                 # foreign key
//       created_at: datetime, default(NOW)
//       updated_at: datetime, default(NOW)
//       title: varchar
//       subtitle: varchar, nullable
//       context: text
//
//     post_images:
//       id: int, primary
//       posts_id: posts(id)                  # another foreign key
//       img_url: varchar
//
//   enums:
//     role:
//     - GENERAL_USER
//     - MODERATOR
//     - ADMIN
//
// The configuration has the following type constraints:
//   primary            sets the type as primary key in the table
//   unique             sets a UNIQUE constraint on the type
//   nullable           undoes the default NOT NULL constraint on a type
//   default(%)         set the default constraint and uses the parameter to set
//                      a default value
//   constraint(%)      sets one or more constraints, constraints can be seperated by ;
//
package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Flags for the application
var (
	migrateFlag = flag.String(`migrate`, ``, `Specifies the output for the SQL migration`)
	dbFlag      = flag.Bool(`db`, false, `Directly inserts the configured structur into the database`)
	dvrFlag     = flag.String(`dvr`, os.Getenv(`DB_DRIVER`), `Sets the driver for the db flag. (postgres by default)`)
	csFlag      = flag.String(`cs`, os.Getenv(`DB_CONNECTION_STRING`), `Sets the connection string for the db flag`)
	modelFlag   = flag.String(`model`, ``, `Sets the output for the Nisevoid/qb models`)
	pkgFlag     = flag.String(`pkg`, ``, `Used by the model flag, sets the package name of the model file(s)`)
)

var mdl Model

func main() {
	flags()

	var errs []error
	success := make(map[string]string)

	if len(mdl.Tables) == 0 {
		flag.PrintDefaults()
		return
	}
	if *modelFlag != `` {
		if file, err := os.Create(*modelFlag); err == nil {
			defer file.Close() // nolint: errcheck

			// nolint: errcheck
			file.WriteString(`// This file is a generated file by github.com/jobstoit/gqb. PLEASE DO NOT EDIT.

package ` + mdl.Pkg + `

import (
	"database/sql"
	"fmt"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/autoqb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
)`)
			CreateQbModel(mdl, file)
			success[`model`] = `written to ` + *modelFlag
		} else {
			errs = append(errs, err)
		}
	}

	if *migrateFlag != `` {
		if file, err := os.Create(*migrateFlag); err == nil {
			defer file.Close() // nolint: errcheck
			CreateMigration(mdl, file)
			success[`migration`] = `written to ` + *migrateFlag
		} else {
			errs = append(errs, err)
		}
	}

	if *dbFlag {
		buff := new(bytes.Buffer)
		CreateMigration(mdl, buff)

		if db, err := sql.Open(mdl.Driver, *csFlag); err == nil {
			defer db.Close() // nolint: errcheck
			if _, err = db.Exec(buff.String()); err != nil {
				errs = append(errs, err)
			}
			success[`db`] = `succesfully executed the migration in the database`
		} else {
			errs = append(errs, err)
		}
	}

	if len(success) > 0 {
		fmt.Printf("\n        \x1b[34m0000      \n     \x1b[34m0000 \x1b[32m000      \x1b[34m##   ##  ##      ##       #####     #####   ##     ##   #####   #####    ##### \n   \x1b[34m0000    \x1b[32m000     \x1b[34m##   ##  ##      ####     ##  ##   ##   ##  ##     ##  ##   ##  ##  ##  ##   ##\n  \x1b[34m0000      \x1b[32m000    \x1b[34m##   ##  ##      ##       #####    #######  ##  #  ##  #######  #####   #######\n \x1b[34m0\x1b[36m000       \x1b[32m0000   \x1b[34m##   ##  ##      ##   ##  ## ##    ##   ##  ##  #  ##  ##   ##  ## ##   ##\n  \x1b[36m000000000 \x1b[32m0000    \x1b[34m######  ######   #####   ##  ##   ##   ##   ### ###   ##   ##  ##  ##   #####\n     \x1b[36m000000000\n\n                               \x1b[32m%s\x1b[0m\n", `succesfully generated`)
		for key, val := range success {
			fmt.Printf("%s: %s\n", key, val)
		}
	}

	for _, err := range errs {
		fmt.Println(err)
	}
}

func flags() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if bitz, err := ioutil.ReadFile(flag.Arg(0)); err == nil {
		mdl = ReadConfig(bitz)
	} else {
		return
	}

	if *dvrFlag != `` {
		mdl.Driver = *dvrFlag
	}

	if mdl.Driver == `` {
		mdl.Driver = `postgres`
	}

	if *pkgFlag != `` {
		mdl.Pkg = *pkgFlag
	}

	if mdl.Pkg == `` {
		mdl.Pkg = `model`
	}
}
