// Copyright 2019 Job Stoit. All rights reserved.

package main

import (
	"os"
)

// Catch is used to panic a function/statement when errors occur
func catch(args ...interface{}) {
	for _, a := range args {
		if err, ok := a.(error); ok && err != nil {
			panic(err)
		}
	}
}

// Create creates a new truncate file
func create(filename string, head ...string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	catch(err)
	for _, h := range head {
		catch(f.Write([]byte(h)))
	}
	return f
}
