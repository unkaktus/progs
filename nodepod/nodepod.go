// nodepod.go - get name of pod on specific node and/or prefix name.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of nodepod, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func init() {
	log.SetFlags(0)
}

func main() {
	var nodename = flag.String("n", "", "node name")
	var podprefix = flag.String("p", "", "pod prefix")
	flag.Parse()

	var template string
	if *nodename == "" {
		template = fmt.Sprintf("{{range .items}}{{.metadata.name}} {{end}}")
	} else {
		template = fmt.Sprintf("{{range .items}}{{if eq .spec.nodeName \"%s\"}}{{.metadata.name}} {{end}}{{end}}", *nodename)
	}
	out, err := exec.Command("kubectl", "get", "pods", "-o=go-template="+template).CombinedOutput()
	if err != nil {
		log.Fatalf("%s", out)
	}
	pods := strings.Split(strings.TrimRight(string(out), "\n "), " ")

	for _, pod := range pods {
		if strings.HasPrefix(pod, *podprefix) {
			fmt.Printf("%s\n", pod)
		}
	}
}
