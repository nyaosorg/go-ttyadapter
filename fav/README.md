go-ttyadapter/fav
=================

`go-ttyadapter/fav` is an alias-like package that selects the preferred backend implementation depending on the environment:

* On Windows with Go 1.20 or earlier, it behaves equivalently to `go-ttyadapter/tty8pe` (using `mattn/go-tty`)
* Otherwise, it behaves equivalently to `go-ttyadapter/tty10pe` (using `golang.org/x/term`)
