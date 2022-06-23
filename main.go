package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort   = "8080"
	messageGoAway = "Go away!"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("Port: %v", port)
	clientType := strings.ToLower(os.Getenv("CLIENT_TYPE"))
	log.Printf("Client: %v", clientType)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, messageGoAway)
	})
	http.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s", strings.TrimPrefix(r.URL.Path, "/s/"))
	})

	var handler func(w http.ResponseWriter, r *http.Request)
	switch clientType {
	case "chrome":
		h, cancel := chromeURL()
		defer cancel()
		handler = h
	case "go":
		handler = netURL()
	case "":
		handler = netURL()
	}

	http.HandleFunc("/url", handler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Couldn't start the server: %s", err)
	}
}
