// kubectx.go - shortcuts for kubectl contexts
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of kubectx, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func ListContexts() error {
	out, err := exec.Command("kubectl", "config", "get-contexts", "-o=name").CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	output := strings.TrimRight(string(out), "\n")
	contexts := strings.Split(output, "\n")
	sort.Strings(contexts)
	cc, err := CurrentContext()
	if err != nil {
		return err
	}
	for _, context := range contexts {
		if context == cc {
			fmt.Print(" ")
		}
		fmt.Printf("%s\n", context)
	}
	return nil
}

func CurrentContext() (string, error) {
	out, err := exec.Command("kubectl", "config", "view", "-o=jsonpath={.current-context}").CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

func RenameContext(from, to string) error {
	out, err := exec.Command("kubectl", "config", "rename-context", from, to).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

func SwitchContext(ctx string) error {
	out, err := exec.Command("kubectl", "config", "use-context", ctx).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
