//go:build !windows

/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import "os"

const eol = "\n"

var initialsPath = os.ExpandEnv("$HOME/.initials")
