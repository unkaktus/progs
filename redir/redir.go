// redir.go - http server to redirect to a link from command line
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of redir, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func LearnEgressAddress() (string, error) {
	tc, err := net.Dial("udp", "1.1.1.1:1")
	if err != nil {
		return "", err
	}
	defer tc.Close()
	host, _, err := net.SplitHostPort(tc.LocalAddr().String())
	if err != nil {
		return "", err
	}
	return host, nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("no link specified")
	}
	link := flag.Args()[0]
	host, err := LearnEgressAddress()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", host+":8888")
	if err != nil {
		l, err = net.Listen("tcp", host+":0")
	}
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, link, http.StatusTemporaryRedirect)
	})
	fmt.Printf("http://%s/\n", l.Addr().String())
	log.Fatal(http.Serve(l, nil))
}
