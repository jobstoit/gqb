// Copyright 2019 Job Stoit. All rights reserved.

// gqb is made for generating a query builder.
//
// The gqb command takes the path to the db.yml as argument.
// So for example:
//   gqb -model ./models/qb.mdl.go db.yml
//
// The command takes the following flags as arguments.
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
//   -go_out            Writes the configuration to NiseVoid/qb model(s) takes
//                      the output file(s) as argument
//
//   -cs_out            Writes te configuration to a dotnet core entityframeworkcore
//                      models.
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
package main

import (
	"flag"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jobstoit/gqb/model"
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

func main() {
	f, err := os.OpenFile(flag.Arg(0))
	if err != nil {
		os.Exit(1)
	}
	mdl := model.FromFile(f)
	if err := f.Close(); err != nil {
		os.Exit(1)
	}

	if modelFlag != `` {
		f, err := os.Create(modelFlag)
		if err == nil {
			template.CreateQbModel(*mdl, f)
		}
	}
}
