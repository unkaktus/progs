// mkremote.go - add git remotes for Go projects
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of mkremote, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.SetFlags(0)
	var remoteName = flag.String("n", "origin", "remote name")
	flag.Parse()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home := os.Getenv("HOME")
		if home == "" {
			log.Fatal("no $HOME set")
		}
		gopath = filepath.Join(home, "go")
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get current directory: %v", err)
	}
	srcPath := filepath.Join(gopath, "src")
	importPath, err := filepath.Rel(srcPath, pwd)
	if err != nil {
		log.Fatalf("unable to extract import path: %v", err)
	}
	remote := "https://" + importPath
	if _, err = url.Parse(remote); err != nil {
		log.Fatalf("produced invalid URL (%s): %v", remote, err)
	}
	log.Printf("adding remote \"%s\" as \"%s\"", remote, *remoteName)
	out, err := exec.Command("git", "remote", "add", *remoteName, remote).CombinedOutput()
	if err != nil {
		log.Fatalf("unable to add remote: %s", out)
	}
}
