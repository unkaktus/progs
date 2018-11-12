// goplay.go - run .go file in Go Playground.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Event struct {
	Delay   int
	Kind    string
	Message string
}

type CompileResponse struct {
	Errors      string
	Events      []Event
	IsTest      bool
	Status      int
	TestsFailed int
}

func playSource(u, source string) (*CompileResponse, error) {
	form := url.Values{}
	form.Add("version", "2")
	form.Add("body", source)
	resp, err := http.PostForm(u+"/compile", form)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("status: %v", resp.Status)
	}
	bb, _ := ioutil.ReadAll(resp.Body)
	cr := &CompileResponse{}
	err = json.Unmarshal(bb, cr)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func run(u, filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	cr, err := playSource(u, string(d))
	if err != nil {
		return err
	}
	if cr.Errors != "" {
		os.Stderr.Write([]byte(cr.Errors))
		return nil
	}
	for _, e := range cr.Events {
		switch e.Kind {
		case "stdout":
			os.Stdout.Write([]byte(e.Message))
		case "stderr":
			os.Stderr.Write([]byte(e.Message))
		}
	}
	return nil
}

func main() {
	var u = flag.String("url", "https://play.golang.org", "url of the playground")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("please provide a go file to play")
	}
	filename := flag.Args()[0]
	if err := run(*u, filename); err != nil {
		log.Fatal(err)
	}
}
