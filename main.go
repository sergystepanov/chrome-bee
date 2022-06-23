package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	h := http.NewServeMux()

	h.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, messageGoAway)
	})
	h.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s", strings.TrimPrefix(r.URL.Path, "/s/"))
	})

	var handler func(w http.ResponseWriter, r *http.Request)
	var cancel []func()
	switch clientType {
	case "chrome":
		handler, cancel = chromeURL()
	case "go":
		handler = netURL()
	case "":
		handler = netURL()
	}

	h.HandleFunc("/url", handler)

	server := &http.Server{Addr: ":" + port, Handler: h}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		log.Printf("Terminate")
		if cancel != nil && len(cancel) > 0 {
			for _, f := range cancel {
				f()
			}
		}

		_ = server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server: %s", err)
	}
}
