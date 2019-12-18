// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"
)

// Catch is used to panic a function/statement when errors occur
func catch(err error, msg string, args ...interface{}) {
	if err != nil {
		fatal(msg, args...)
	}
}

// Fatal closes the program with a message
func fatal(msg string, args ...interface{}) {
	if !strings.HasPrefix(msg, "\n") {
		msg += "\n"
	}

	fmt.Printf(msg, args...)
	os.Exit(1)
}
