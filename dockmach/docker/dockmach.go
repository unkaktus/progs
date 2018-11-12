// dockmach.go - run docker via docker-machine across terminals without handwork.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

// Rename real `docker` binary to `docker.real`, install this as `docker`:
// $ go get -v github.com/nogoegst/progs/dockmach/docker
// Pick your current Docker Machine:
// $ echo 'kitchen-sink' > ~/.docker/machine/current
// Run docker from any terminal,
// Have fun.

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if os.Args[0] != "docker" {
		log.Fatal("I wasn't called as `docker`")
	}
	machineBin, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".docker/machine/current"))
	if err != nil {
		log.Fatal(err)
	}
	machine := strings.TrimSuffix(string(machineBin), "\n")

	out, err := exec.Command("docker-machine", "env", machine).CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	env := []string{}
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		envVar := strings.Replace(strings.TrimPrefix(text, "export "), "\"", "", -1)
		env = append(env, envVar)
	}
	cmd := exec.Command("docker.real", os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), env...)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
