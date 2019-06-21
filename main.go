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

import "fmt"

func main() {
	fmt.Println("vim-go")
}
