package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func ReleaseNameByChartName(chartName string, namespace string) (string, error) {
	args := []string{"list", "--deployed"}
	if namespace != "" {
		args = append(args, []string{"--namespace", namespace}...)
	}

	out, err := exec.Command("helm", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", out)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		l := scanner.Text()
		fields := strings.FieldsFunc(l, func(r rune) bool {
			return r == '\t'
		})
		for i, f := range fields {
			fields[i] = strings.Trim(f, " ")
		}
		releaseName := fields[0]
		chart := fields[4]
		if strings.HasPrefix(chart, chartName+"-") {
			return releaseName, nil
		}
	}
	return "", errors.New("no release of this chart")
}

func SetReleaseImageTag(chartName, releaseName, imageTag string, wait bool) error {
	args := []string{"upgrade", "--reuse-values", "--set", "image.tag=" + imageTag, releaseName, chartName}
	if wait {
		args = append(args, "--wait")
	}
	out, err := exec.Command("helm", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}
	log.Printf("%s", out)
	return nil
}

func main() {
	log.SetFlags(0)
	var wait = flag.Bool("wait", false, "wait release to finish its rollout")
	var namespace = flag.String("namespace", "", "namespace to look releases in")
	flag.Parse()
	switch len(flag.Args()) {
	case 0:
		log.Fatal("no chart name specified")
	case 1:
		log.Fatal("no image tag specified")
	case 2:
		break
	default:
		log.Fatal("too many arguments")
	}
	chartName := flag.Args()[0]
	imageTag := flag.Args()[1]
	releaseName, err := ReleaseNameByChartName(chartName, *namespace)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("upgrading '%s'", releaseName)
	if err := SetReleaseImageTag(chartName, releaseName, imageTag, *wait); err != nil {
		log.Fatal(err)
	}
}
