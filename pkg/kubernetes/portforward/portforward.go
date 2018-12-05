package portforward

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

type PortForwarder interface {
	Addr() string
	Close()
}

func getForwardAddress(r io.Reader) (string, error) {
	s, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return "", err
	}
	s = strings.TrimPrefix(s, "Forwarding from ")
	sp := strings.Split(s, " -> ")
	return sp[0], nil
}

type portForwarder struct {
	addr   string
	cancel context.CancelFunc
}

func (pfw *portForwarder) Close() {
	pfw.cancel()
}

func (pfw *portForwarder) Addr() string {
	return pfw.addr
}

func NewPortForward(ctx context.Context, podname, port string, opts ...string) (*portForwarder, error) {
	ctx, cancel := context.WithCancel(ctx)
	kubectlArgs := append([]string{"port-forward", podname, ":" + port}, opts...)

	cmd := exec.CommandContext(ctx, "kubectl", kubectlArgs...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err

	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go cmd.Run()

	addr, err := getForwardAddress(stdoutPipe)
	if err != nil {
		stderr, _ := ioutil.ReadAll(stderrPipe)
		return nil, errors.New(string(stderr))
	}

	pfw := &portForwarder{
		addr:   addr,
		cancel: cancel,
	}
	return pfw, nil
}
