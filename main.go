/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"runtime"

	"github.com/electrikmilk/ttuy"
)

var EOL = "\n"

func main() {
	if runtime.GOOS == "windows" {
		EOL = "\r\n"
	}
	registerArg("help", "h", "Show this help message")
	registerArg("remote", "r", "Use remote branches as basis")
	registerArg("initials", "i", "Set new initials")
	registerArg("debug", "d", "Get debug output on error")
	if arg("help") {
		usage()
		return
	}
	checkForGit()
	if arg("initials") {
		saveInitials()
	} else {
		getInitials()
	}
	var branchType string
	ttuy.Menu("Type of Branch", []ttuy.Option{
		{
			Label: "Hotfix",
			Callback: func() {
				branchType = "hotfix"
			},
		},
		{
			Label: "Bug",
			Callback: func() {
				branchType = "bug"
			},
		},
		{
			Label: "Enhancement",
			Callback: func() {
				branchType = "enhancement"
			},
		},
		{
			Label: "Feature",
			Callback: func() {
				branchType = "feature"
			},
		},
	})
	var reference string
	fmt.Println(ttuy.Style("Enter a ticket number, or dash seperated string describing the branch.", ttuy.Dim))
	ttuy.Ask("Reference", &reference)
	var name string = fmt.Sprintf("%s/%s/%s", branchType, initials, reference)
	create(&name)
}

// Handle a program error
func handle(label string, err error) {
	if err != nil {
		if arg("debug") {
			fmt.Println(label)
			panic(err)
		} else {
			ttuy.FailErr(label, err)
		}
	}
}

// Alias for handle, defaults label to "Git"
func handleGit(err error) {
	handle("Git", err)
}
