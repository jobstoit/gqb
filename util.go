// Copyright 2019 Job Stoit. All rights reserved.

package main

// Catch is used to panic a function/statement when errors occur
func catch(args ...interface{}) {
	for _, a := range args {
		if err, ok := a.(error); ok && err != nil {
			panic(err)
		}
	}
}
