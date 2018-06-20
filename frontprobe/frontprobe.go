// frontprobe.go - quickly test domain fronting availability.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(0)
	var dump = flag.Bool("dump", false, "dump response body")
	var hostname = flag.String("a", "", "TCP address")
	var front = flag.String("n", "", "SNI")
	var hostHeader = flag.String("h", "", "Host header")
	var path = flag.String("p", "/", "path")
	flag.Parse()

	req, err := http.NewRequest(http.MethodGet, "https://"+*hostname+*path, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Host = *hostHeader
	t := http.DefaultTransport.(*http.Transport)
	t.TLSClientConfig = &tls.Config{
		ServerName:         *front,
		InsecureSkipVerify: true,
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("status code: %v", resp.StatusCode)
	location := resp.Header.Get("Location")
	if location != "" {
		log.Printf("location: %s", location)
	}
	if *dump {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("body: %s", body)
	}
}
