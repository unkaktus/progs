progs are tiny utilities written in Go.

but why?

 * readablity (no more `if [ $# -eq 42]; then echo $1; fi;` 
 * static binary (contains everything, no need to bring any userspace stuff)
 * painless crosscompilation (just set GOOS, GOARCH and run `go build`)
 * extensibility (there are tons of libraries one can use)
 * ???
 * the gopher!!
