// chrome_tmp.go - open chrome in incognito mode with ephemeral profile
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetFlags(0)
	args := os.Args[1:]

	dir, err := ioutil.TempDir("", "chrome-tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	cmd := exec.Command("chrome", append([]string{"--incognito", "--user-data-dir=" + dir}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
