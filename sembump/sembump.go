// sembump.go - bump semantic versions in git
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of sembump, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/masterminds/semver"
)

func GitCurrentTag() (string, error) {
	out, err := exec.Command("git", "describe", "--tags").CombinedOutput()
	if err != nil {
		gitNoTagsError := "No names found"
		if strings.Contains(string(out), gitNoTagsError) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

func GitTag(tag string) error {
	out, err := exec.Command("git", "tag", tag).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}
	return nil
}

func main() {
	log.SetFlags(0)

	var majorFlag = flag.Bool("M", false, "major bump")
	var minorFlag = flag.Bool("m", false, "minor bump")
	var patchFlag = flag.Bool("p", false, "patch bump")
	flag.Parse()

	flagCount := 0
	if *majorFlag {
		flagCount++
	}
	if *minorFlag {
		flagCount++
	}
	if *patchFlag {
		flagCount++
	}

	if flagCount != 1 {
		log.Fatal("only one flag can be specified")
	}

	tag, err := GitCurrentTag()
	if err != nil {
		log.Fatalf("unable to get git tag: %v", err)
	}
	if tag == "" {
		tag = "v0.0.0"
	}

	version, err := semver.NewVersion(tag)
	if err != nil {
		log.Fatalf("unable to parse semver: %v", err)
	}

	var v semver.Version
	switch {
	case *majorFlag:
		v = version.IncMajor()
	case *minorFlag:
		v = version.IncMinor()
	case *patchFlag:
		v = version.IncPatch()
	}
	tag = "v" + v.String()
	log.Printf("tagging with %s", tag)
	if err := GitTag(tag); err != nil {
		log.Fatalf("unable to git tag: %v", err)
	}
}
