// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
// Package main wires all of the middleware and services together
package main

import (
	"fmt"
	"os"

	"github.ibm.com/Alchemy-Key-Protect/key-management-api/cmd"
)

//  semver and commit are set by build for runtime environments
var semver string
var commit string
var runtime string // not yet used

func main() {
	cmd.SetVersion(semver, commit)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
