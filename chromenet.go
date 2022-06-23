package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func chromeURL() (func(w http.ResponseWriter, r *http.Request), context.CancelFunc) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)

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
			log.Fatal(err)
		}

		charset := findCharset(buf)
		log.Printf("found charset: %v", charset)

		w.Header().Set("Content-Type", "text/html; charset="+charset)
		_, _ = fmt.Fprintf(w, "%s", buf)
	}, cancel
}

func getHTML(url string, body *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", body, chromedp.ByQuery),
	}
}
