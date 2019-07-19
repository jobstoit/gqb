// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"fmt"
	"os"
)

// Catch is used to panic a function/statement when errors occur
func catch(err error, msg string, args ...interface{}) {
	if err != nil {
		fatal(msg, args...)
	}
}

// Fatal closes the program with a message
func fatal(msg string, args ...interface{}) {
	if len(args) == 0 {
		msg += "\n"
	}

	fmt.Printf(msg, args...)
	os.Exit(1)
}
