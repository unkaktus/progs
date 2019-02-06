package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/nogoegst/progs/pkg/egressip"
	"github.com/nogoegst/textqr"
	"github.com/nogoegst/tlspin"
)

func main() {
	host, err := egressip.Get()
	if err != nil {
		log.Fatal(err)
	}
	l, err := tlspin.Listen("tcp", "0.0.0.0:0", "whateverkey")
	if err != nil {
		log.Fatal(err)
	}
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	addr := net.JoinHostPort(host, port)

	log.Printf("serving at https://%s", addr)
	_, err = textqr.Write(os.Stdout, "https://"+addr, textqr.L, true, false)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.String()
		u = strings.Replace(u, "/https:/", "https://", -1)
		u = strings.Replace(u, "/http:/", "http://", -1)
		fmt.Printf("here: %s\n", u)
	})
	if err := http.Serve(l, nil); err != nil {
		log.Fatal(err)
	}
}
