# Arbor

Create git branches using a template:

```
type/initials/reference
```

![Example](https://i.imgur.com/I4nyxIY.gif)

You're prompted for your initials on first run, but you can set new initials with the `-i` or `--initials` flag or by editing the `.initials` file in your home directory.

You can base your branch off of remote branches using the `-r` or `--remote` flag.

Arbor does not currently support remote protected branches.

Arbor uses my TUI framework [ttuy](https://github.com/electrikmilk/ttuy).
