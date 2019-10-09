package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {

	path := os.Args[1]
	revision := os.Args[2]

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		panic(err)
	}

	fmt.Println(h.String())
}
