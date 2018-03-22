package main

import (
	"encoding/base64"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const (
	CAkey = "/etc/kubernetes/pki/ca.key"
	CAcrt = "/etc/kubernetes/pki/ca.crt"
)

func main() {
	log.SetFlags(0)
	var days = flag.String("days", "500", "days of validity")
	flag.Parse()
	csrB64 := flag.Args()[0]
	csrData, err := base64.StdEncoding.DecodeString(csrB64)
	if err != nil {
		log.Fatal(err)
	}
	csr, err := ioutil.TempFile("/tmp", "kubecred")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(csr.Name())
	err = pem.Encode(csr, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrData,
	})
	if err != nil {
		log.Fatal(err)
	}
	csr.Close()

	crt, err := ioutil.TempFile("/tmp", "kubecred")
	if err != nil {
		log.Fatal(err)
	}
	crt.Close()
	defer os.Remove(crt.Name())
	out, err := exec.Command("openssl", "x509", "-req", "-in", csr.Name(), "-CA", CAcrt, "-CAkey", CAkey, "-CAcreateserial", "-out", crt.Name(), "-days", *days).CombinedOutput()
	if err != nil {
		log.Fatalf("%s", out)
	}
	crtPEM, err := ioutil.ReadFile(crt.Name())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", base64.StdEncoding.EncodeToString(crtPEM))
}
