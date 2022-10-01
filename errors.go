/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"github.com/electrikmilk/ttuy"
)

// Handle a program error
func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// Gracefully handle a git related error, as it's sometimes not an application error
func handleGit(err error) {
	ttuy.FailErr("Git", err)
}
