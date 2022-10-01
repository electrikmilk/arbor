/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/electrikmilk/ttuy"
)

var initials string
var initialsPath string = os.ExpandEnv("$HOME/.initials")

func getInitials() {
	if _, err := os.Stat(initialsPath); errors.Is(err, os.ErrNotExist) || arg("initials") == true {
		saveInitials()
		fmt.Printf(ttuy.Style("Use -i flag to reset initials."+EOL+EOL, ttuy.Dim))
	} else {
		savedInitials, readErr := os.ReadFile(initialsPath)
		handle(readErr)
		initials = string(savedInitials)
	}
}

func saveInitials() {
	if _, err := os.Stat(initialsPath); errors.Is(err, os.ErrNotExist) {
		f, createErr := os.Create(initialsPath)
		ttuy.FailErr("Unable to create initials", createErr)
		defer func(f *os.File) {
			err := f.Close()
			handle(err)
		}(f)
	}
	ttuy.Ask("Enter your initials", &initials)
	writeErr := os.WriteFile(initialsPath, []byte(initials), 0774)
	ttuy.FailErr("Unable to save initials", writeErr)
	ttuy.Success("Initials saved!")
}
