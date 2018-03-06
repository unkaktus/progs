// kubewrap.go - make cli tools connect to Kubernetes pods directly.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of kubewrap, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nogoegst/progs/pkg/kubernetes/portforward"
)

const DefaultPort = "80"

func listPods() ([]string, error) {
	template := "{{range .items}}{{.metadata.name}} {{end}}"
	out, err := exec.Command("kubectl", "get", "pods", "-o=go-template="+template).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s", out)
	}
	pods := strings.Split(strings.TrimRight(string(out), "\n "), " ")
	return pods, nil
}

func containsPod(pods []string, s string) (string, bool) {
	for _, pod := range pods {
		if strings.Contains(s, pod) {
			return pod, true
		}
	}
	return "", false
}

func extractAddress(s, pod string) string {
	addr := pod
	rest := strings.SplitN(s, pod, 2)[1]
	if strings.HasPrefix(rest, ":") {
		addr += extractPortSuffix(rest)
	}
	return addr
}

func extractPortSuffix(s string) string {
	for i := 0; i < len(s)-1; i++ {
		p := s[1 : 2+i]
		_, e := strconv.Atoi(p)
		if e != nil {
			return s[:1+i]
		}
		if i == len(s)-2 {
			return s[:2+i]
		}
	}
	return ""
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("specify app name")
	}
	app := os.Args[1]
	args := os.Args[2:]

	if app != "kubectl" {
		var pfw portforward.PortForwarder
		ctx := context.TODO()

		pods, err := listPods()
		if err != nil {
			log.Fatal(err)
		}

		for i, arg := range args {
			pod, ok := containsPod(pods, arg)
			if !ok {
				continue
			}

			addr := extractAddress(arg, pod)
			_, port, _ := net.SplitHostPort(addr)
			if port == "" {
				port = DefaultPort
			}

			if pfw == nil {
				pfw, err = portforward.NewPortForward(ctx, pod, port)
				if err != nil {
					log.Fatal(err)
				}
				defer pfw.Close()
			}
			args[i] = strings.Replace(arg, addr, pfw.Addr(), 1)
		}
	}
	cmd := exec.Command(app, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
