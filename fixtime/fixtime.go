// fixtime.go - crudely override system clock with RTC clock value.
// This is only useful for system clocks which are not ticking,
// so hwclock doesn't work. :\
// The precision if this thing is 1s as it reads RTC data from since_epoch.
//
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

// +build linux

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

func readRTC() (int64, error) {
	b, err := ioutil.ReadFile("/sys/class/rtc/rtc0/since_epoch")
	if err != nil {
		return 0, err
	}
	str := strings.TrimRight(string(b), "\n")
	s, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return s, nil
}

func main() {
	var intervalFlag = flag.String("i", "1s", "interval")
	flag.Parse()
	interval, err := time.ParseDuration(*intervalFlag)
	if err != nil {
		log.Fatal(err)
	}
	for {
		s, err := readRTC()
		if err != nil {
			log.Fatal(err)
		}
		err = unix.Settimeofday(&unix.Timeval{Sec: s})
		if err != nil {
			log.Fatal(err)
		}
		print("⌚ ")

		time.Sleep(interval)
	}
}
