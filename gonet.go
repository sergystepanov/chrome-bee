package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const defaultUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"

func netURL() func(w http.ResponseWriter, r *http.Request) {
	ct := http.DefaultTransport.(*http.Transport).Clone()
	ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var client = &http.Client{Transport: ct}

	ua := os.Getenv("UA")
	if ua == "" {
		ua = defaultUA
	}
	log.Printf("UA: %v", ua)

	return func(w http.ResponseWriter, r *http.Request) {
		// parse POST params
		if err := r.ParseForm(); err != nil {
			log.Printf("couldn't parse params, %v", err)
			return
		}

		url := r.Form.Get("_url")
		if url == "" {
			log.Printf("bad params!")
			return
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("bad url: %v", url)
			return
		}

		req.Header.Set("User-Agent", ua)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("couldn't get, %v", err)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("has close error, %v", err)
				return
			}
		}()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("couldn't read, %v", err)
			return
		}

		charset := findCharset(string(body))
		log.Printf("found charset: %v", charset)

		w.Header().Set("Content-Type", "text/html; charset="+charset)
		_, _ = fmt.Fprintf(w, "%s", body)
	}
}
