package main

import (
	"fmt"
	"log"

	"github.com/nogoegst/progs/pkg/egressip"
)

func main() {
	a, err := egressip.Get()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}
