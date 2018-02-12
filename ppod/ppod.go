// ppod.go - fetch CPU profile and trace from running k8s pod.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of pod, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func FetchProfile(w io.Writer, addr string, d time.Duration) error {
	a := fmt.Sprintf("http://%s/debug/pprof/profile?seconds=%d", addr, int(d.Seconds()))
	resp, err := http.Get(a)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

func FetchTrace(w io.Writer, addr string, d time.Duration) error {
	a := fmt.Sprintf("http://%s/debug/pprof/trace?seconds=%d", addr, int(d.Seconds()))
	resp, err := http.Get(a)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

func GetForwardAddress(r io.Reader) (string, error) {
	s, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return "", err
	}
	s = strings.TrimPrefix(s, "Forwarding from ")
	sp := strings.Split(s, " -> ")
	return sp[0], nil
}

func main() {
	log.SetFlags(0)
	var port = flag.String("p", "6060", "pod port")
	var profile = flag.String("profile", "", "profile duration (30s default)")
	var trace = flag.String("trace", "", "trace duration (30s default)")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("no pod specified")
	}

	podname := flag.Args()[0]

	cmd := exec.Command("kubectl", "port-forward", podname, ":"+*port)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go cmd.Run()

	addr, err := GetForwardAddress(stdoutPipe)
	if err != nil {
		stderr, _ := ioutil.ReadAll(stderrPipe)
		log.Fatalf("%s", stderr)
	}

	if *profile == "" && *trace == "" {
		*profile = "30s"
		*trace = "30s"
	}

	if *profile != "" {
		profileDuration, err := time.ParseDuration(*profile)
		if err != nil {
			log.Fatal(err)
		}
		profileFile, err := os.Create(podname + ".profile")
		if err != nil {
			log.Fatal(err)
		}
		defer profileFile.Close()
		log.Printf("fetching profile to %s", profileFile.Name())
		err = FetchProfile(profileFile, addr, profileDuration)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *trace != "" {
		traceDuration, err := time.ParseDuration(*trace)
		if err != nil {
			log.Fatal(err)
		}
		traceFile, err := os.Create(podname + ".trace")
		if err != nil {
			log.Fatal(err)
		}
		defer traceFile.Close()
		log.Printf("fetching trace to %s", traceFile.Name())
		err = FetchTrace(traceFile, addr, traceDuration)
		if err != nil {
			log.Fatal(err)
		}
	}
}
