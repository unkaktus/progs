// bundle.go - git tag go build docker build docker push tool
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of bundle, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GitCurrentTag() (string, error) {
	out, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

func GoBuild(tempDir string) error {
	binaryPath := filepath.Join(tempDir, "app")
	cmd := exec.Command("go", "build", "-v", "-o", binaryPath)
	goos := "linux"
	if o := os.Getenv("GOOS"); o != "" {
		goos = o
	}
	cmd.Env = append(os.Environ(), "GOOS="+goos)
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}

const MinimalDockerfileFormat = "FROM %s\nCOPY app /\nCMD [\"/app\"]"

func WriteMinimalDockerfile(tmpDir, baseImage string) error {
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	MinimalDockerfile := fmt.Sprintf(MinimalDockerfileFormat, baseImage)
	return ioutil.WriteFile(dockerfilePath, []byte(MinimalDockerfile), 0400)
}

func DockerBuild(image, tag string, tmpDir string) error {
	imageTag := image + ":" + tag
	cmd := exec.Command("docker", "build", "-t", imageTag, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DockerPush(image, tag string) error {
	imageTag := image + ":" + tag
	cmd := exec.Command("docker", "push", imageTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var baseImage = flag.String("b", "scratch", "image name")
	var imageName = flag.String("n", "", "image name")
	var tagFlag = flag.String("t", "", "image tag")
	var push = flag.Bool("p", false, "push image")
	flag.Parse()
	if *baseImage == "" {
		log.Fatalf("no image name specified")
	}
	var tag string
	if *tagFlag == "" {
		var err error
		tag, err = GitCurrentTag()
		if err != nil {
			log.Fatalf("unable to get git current tag: %v", err)
		}
	} else {
		tag = *tagFlag
	}
	log.Printf("using tag: %v", tag)

	tempDir, err := ioutil.TempDir("", "bundle")
	if err != nil {
		log.Fatalf("unable to create temp dir: %v", err)
	}
	log.Printf("temporary directory is %s", tempDir)
	defer os.RemoveAll(tempDir)

	if err := GoBuild(tempDir); err != nil {
		log.Fatalf("go build failed: %v", err)
	}

	if err := WriteMinimalDockerfile(tempDir, *baseImage); err != nil {
		log.Fatalf("unable to write Dockerfile: %v", err)
	}

	if err := DockerBuild(*imageName, tag, tempDir); err != nil {
		log.Fatalf("docker build failed: %v", err)
	}

	if *push {
		if err := DockerPush(*imageName, tag); err != nil {
			log.Fatal("docker push failed: %v", err)
		}
	}
	log.Printf("done!")
}
