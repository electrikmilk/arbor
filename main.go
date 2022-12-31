/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"

	"github.com/electrikmilk/ttuy"
)

func main() {
	customUsage = "[branch|commit]"
	registerArg("help", "h", "Show this help message")
	registerArg("remote", "r", "Use remote branches as basis")
	registerArg("initials", "i", "Set new initials")
	registerArg("debug", "d", "Get debug output on error")
	if arg("help") || len(os.Args) <= 1 {
		usage()
		return
	}
	checkForGit()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "branch":
			startBranch()
		case "commit":
			createCommit()
		}
	}
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
