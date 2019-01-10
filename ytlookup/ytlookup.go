// ytlookup.go - get top hit YouTube search results as list of videoIDs.
// E.g. stream "Roisin Murphy - Ten Miles High":
// youtube-dl -o - $(ytlookup roisin murphy ten miles high) | ffplay -
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of progs, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/savaki/jq"
)

func main() {
	var n = flag.Int("n", 1, "max number of ids to return")
	flag.Parse()
	query := strings.Join(flag.Args(), " ")

	u, _ := url.Parse("https://www.youtube.com/results?pbj=1")
	v := u.Query()
	v.Set("search_query", query)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("x-youtube-client-version", "2.20190109")
	req.Header.Set("x-youtube-client-name", "1")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	op, _ := jq.Parse(".[1].response.contents.twoColumnSearchResultsRenderer.primaryContents.sectionListRenderer.contents.[0].itemSectionRenderer.contents")
	value, err := op.Apply(body)
	if err != nil {
		log.Fatal(err)
	}

	contents := []struct {
		VideoRenderer struct {
			VideoID string `json:"videoId"`
		} `json:"videoRenderer"`
	}{}
	err = json.Unmarshal(value, &contents)
	if err != nil {
		log.Fatal(err)
	}
	videoIDs := []string{}
	for _, c := range contents {
		if c.VideoRenderer.VideoID != "" {
			videoIDs = append(videoIDs, c.VideoRenderer.VideoID)
		}
	}

	for i, videoID := range videoIDs {
		if i+1 > *n {
			break
		}
		fmt.Printf("%s\n", videoID)
	}
}
