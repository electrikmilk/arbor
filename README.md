# Arbor

[![Releases](https://img.shields.io/github/v/release/electrikmilk/arbor?include_prereleases)](https://github.com/electrikmilk/arbor/releases)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/electrikmilk/arbor?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrikmilk/arbor)](https://goreportcard.com/report/github.com/electrikmilk/arbor)

## Current features

### Auto commit

```console
arbor commit
```

This mimics a GitHub feature but for the command line. If you have only one file in staging, this will automatically commit it with an appropriate imparative prefix (e.g. Add file.txt).

### Branch template

Create git branches using a standard template:

```
type/initials/reference
```

```console
arbor branch
```

![Example](https://i.imgur.com/I4nyxIY.gif)

You're prompted for your initials on first run, but you can set new initials with the `-i` or `--initials` flag or by editing the `.initials` file in your home directory.

You can base your branch off of remote branches using the `-r` or `--remote` flag.

Arbor does not currently support remote protected branches.

Arbor uses my TUI framework [ttuy](https://github.com/electrikmilk/ttuy).
