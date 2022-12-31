/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"

	"github.com/electrikmilk/args-parser"
	"github.com/electrikmilk/ttuy"
)

func main() {
	args.CustomUsage = "[branch|commit]"
	args.Register("help", "h", "Show this help message")
	args.Register("remote", "r", "Use remote branches as basis")
	args.Register("initials", "i", "Set new initials")
	args.Register("debug", "d", "Get debug output on error")
	if args.Using("help") || len(os.Args) <= 1 {
		args.PrintUsage()
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
		if args.Using("debug") {
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
