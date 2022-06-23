package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chromedp/chromedp"
)

func chromeURL() (func(w http.ResponseWriter, r *http.Request), []func()) {
	ua := os.Getenv("UA")
	if ua == "" {
		ua = defaultUA
	}
	log.Printf("UA: %v", ua)

	var cleanup []func()

	dir, err := ioutil.TempDir("", "chrome")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Chrome dir: %v", dir)
	cleanup = append(cleanup, func() {
		_ = os.RemoveAll(dir)
	})

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.UserAgent(ua),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	cleanup = append(cleanup, cancel)

	// also set up a custom logger
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	cleanup = append(cleanup, cancel)

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

		// capture screenshot of an element
		var buf string
		// capture entire browser viewport, returning png with quality=90
		if err := chromedp.Run(ctx, getHTML(url, &buf)); err != nil {
			_, _ = fmt.Fprintf(w, "%s", err)
			return
		}

		charset := findCharset(buf)
		log.Printf("found charset: %v", charset)

		w.Header().Set("Content-Type", "text/html; charset="+charset)
		_, _ = fmt.Fprintf(w, "%s", buf)
	}, cleanup
}

func getHTML(url string, body *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.WaitNotPresent(),
		chromedp.OuterHTML("html", body, chromedp.ByQuery),
	}
}
