// main.go - kubectx cli
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of kubectx, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"errors"
	"flag"
	"log"
	"strings"
)

func Usage() {
	log.Println("Usage of kubectx:")
	log.Println("\t$ kubectx\t\tshow contexts with the current one shifted")
	log.Println("\t$ kubectx ctx\t\tswitch to context \"ctx\"")
	log.Println("\t$ kubectx ctx2=ctx1\t\trename context \"ctx1\" to \"ctx2\"")
}

var InvalidCommandErr = errors.New("invalid command")

func main() {
	log.SetFlags(0)
	var helpFlag = flag.Bool("h", false, "show help")
	flag.Parse()
	if *helpFlag {
		Usage()
		return
	}

	var err error
	switch len(flag.Args()) {
	case 0:
		err = ListContexts()
	case 1:
		contexts := strings.Split(flag.Args()[0], "=")
		switch len(contexts) {
		case 2:
			err = RenameContext(contexts[1], contexts[0])
		case 1:
			err = SwitchContext(contexts[0])
		default:
			err = InvalidCommandErr
		}
	default:
		err = InvalidCommandErr
	}

	if err == InvalidCommandErr {
		Usage()
	}

	if err != nil {
		log.Fatalf("kubectx: %v", err)
	}
}
