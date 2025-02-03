package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var serviceMap = map[string]string{
	"test1":        "http://localhost:3450",
	"user-service": "http://localhost:8082",
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Get("Authorization")
	if cookie == "" {
		log.Println("Authorization not found")
	} else {
		log.Println("Authorization:", cookie)
	}

	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	if len(parts) < 2 {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	serviceName, restPath := parts[0], parts[1]
	targetService, exists := serviceMap[serviceName]
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	targetURL, err := url.Parse(targetService)
	if err != nil {
		http.Error(w, "Failed to parse target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Настраиваем поведение прокси
	proxy.Director = func(req *http.Request) {
		req.Header = r.Header
		req.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = "/" + restPath
		log.Printf("Proxying to: %s%s", targetService, req.URL.Path)
	}

	// Передаем запрос в прокси
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", proxyHandler)

	port := 8080
	fmt.Printf("API Gateway running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
