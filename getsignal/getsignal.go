// getsinal.go - fetch latest Signal APK.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const latestURL = `https://updates.signal.org/android/latest.json`

type LatestDescription struct {
	SHA256Sum   string `json:"sha256sum"`
	URL         string `json:"url"`
	VersionCode int    `json:"versionCode"`
	VersionName string `json:"versionName""`
}

func fetchLatestDescription(u string) (*LatestDescription, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	latest := &LatestDescription{}
	err = json.Unmarshal(body, latest)
	if err != nil {
		return nil, err
	}
	return latest, nil
}

func fetchSignal(d *LatestDescription) ([]byte, error) {
	resp, err := http.Get(d.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	h := sha256.New()
	buf := bytes.Buffer{}
	_, err = io.Copy(io.MultiWriter(h, &buf), resp.Body)
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(h.Sum(nil))
	if hashString != d.SHA256Sum {
		return nil, fmt.Errorf("hash mismatch: want %s, got %s", d.SHA256Sum, hashString)
	}
	return buf.Bytes(), nil
}

func main() {
	log.SetFlags(0)
	latest, err := fetchLatestDescription(latestURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Fetching Signal APK v%s", latest.VersionName)
	apkData, err := fetchSignal(latest)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("Signal-v"+latest.VersionName+".apk", apkData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
