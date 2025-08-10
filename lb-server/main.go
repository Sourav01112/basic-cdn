package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	edgeServers := os.Getenv("EDGE_SERVERS")
	if edgeServers == "" {
		log.Fatal("EDGE_SERVERS environment variable required")
	}

	serverURL, err := url.Parse("http://" + edgeServers)
	if err != nil {
		log.Fatal("Invalid edge server URL:", err)
	}

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(serverURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Load Balancer: Routing %s to %s", r.URL.Path, serverURL.Host)

		r.Header.Set("X-CDN-Server", "load-balancer")

		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Load Balancer OK"))
	})

	log.Println("Load Balancer starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
