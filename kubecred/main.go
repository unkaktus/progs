package main

import (
        "bufio"
        "crypto/rand"
        "crypto/rsa"
        "crypto/x509"
        "crypto/x509/pkix"
        "encoding/base64"
        "encoding/pem"
        "flag"
        "io/ioutil"
        "log"
        "os"
        "os/exec"
        "path/filepath"
)

func main() {
        log.SetFlags(0)
        flag.Parse()
        keyName := flag.Args()[0]
        cn := flag.Args()[1]
        sk, err := rsa.GenerateKey(rand.Reader, 2048)
        if err != nil {
                log.Fatal(err)
        }
        skPEM := pem.EncodeToMemory(&pem.Block{
                Type:  "RSA PRIVATE KEY",
                Bytes: x509.MarshalPKCS1PrivateKey(sk),
        })
        if err := ioutil.WriteFile(keyName+".key", skPEM, 0644); err != nil {
                log.Fatal(err)
        }

        subj := pkix.Name{
                CommonName: cn,
        }
        csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{Subject: subj}, sk)
        if err != nil {
                log.Fatal(err)
        }
        csrB64 := base64.StdEncoding.EncodeToString(csr)

        log.Printf("Copy CSR to kubecred-ca:\n%s", csrB64)
        log.Printf("Enter response from kubecred-ca:\n")
        crtB64, err := bufio.NewReader(os.Stdin).ReadString(byte('\n'))
        if err != nil {
                log.Fatal(err)
        }
        crtData, err := base64.StdEncoding.DecodeString(crtB64)
        if err != nil {
                log.Fatal(err)
        }
        err = ioutil.WriteFile(keyName+".crt", crtData, 0644)
        if err != nil {
                log.Fatal(err)
        }
        log.Printf("Result written to %s.key and %s.crt", keyName, keyName)
        keyAbs, _ := filepath.Abs(keyName + ".key")
        crtAbs, _ := filepath.Abs(keyName + ".crt")
        out, err := exec.Command("kubectl", "config", "set-credentials", keyName, "--client-certificate", crtAbs, "--client-key", keyAbs).CombinedOutput()
        if err != nil {
                log.Fatalf("%s", out)
        }
        log.Printf("%s", out)

}
