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
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	addr, err := GetForwardAddress(stdoutPipe)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("addr: %s", addr)

	profileFile, err := os.Create(podname + ".profile")
	if err != nil {
		log.Fatal(err)
	}
	defer profileFile.Close()
	log.Print("fetching profile...")
	err = FetchProfile(profileFile, addr, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	traceFile, err := os.Create(podname + ".trace")
	if err != nil {
		log.Fatal(err)
	}
	defer traceFile.Close()
	log.Print("fetching trace...")
	err = FetchTrace(traceFile, addr, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

}
