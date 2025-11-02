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
	args.Register(args.Argument{
		Name:        "help",
		Short:       "h",
		Description: "Show this help message",
	})
	args.Register(args.Argument{
		Name:        "remote",
		Short:       "r",
		Description: "Use remote branches as basis",
	})
	args.Register(args.Argument{
		Name:        "initials",
		Short:       "i",
		Description: "Set new initials",
	})
	args.Register(args.Argument{
		Name:        "debug",
		Short:       "d",
		Description: "Get debug output on error",
	})

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
