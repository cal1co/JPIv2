package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/cal1co/jpiv2/dboperations"
	"github.com/cal1co/jpiv2/handlers"
	opensearch "github.com/opensearch-project/opensearch-go"
	"golang.org/x/time/rate"
)

const IndexName = "go-test-index1"

var limiter = NewIPRateLimiter(1, 5)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func main() {
	// Initialize the client with SSL/TLS enabled.
	client := initClient()

	// Print OpenSearch version information on console.
	fmt.Println(client.Info())

	// Create an index with non-default settings.
	res := dboperations.CreateIndex(IndexName)
	fmt.Println("Creating index")
	fmt.Println(res)

	// Insert
	// dictconfig.AddEntries(IndexName, client)
	// Define the API endpoints.

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("OK")
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.SearchHandler(w, r, client, IndexName)
	})

	mux.HandleFunc("/generatekey", handlers.GenerateKeyHandler)

	// Start the HTTP server.
	fmt.Println("Server listening on port 8888")

	if err := http.ListenAndServe(":8888", limitMiddleware(mux)); err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	// dboperations.Cleanup(IndexName, client)
}

func initClient() *opensearch.Client {
	fmt.Println("INITIATING CLIENT")
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses:     []string{"http://localhost:9200"},
		MaxRetries:    5,
		RetryOnStatus: []int{502, 503, 504},
	})
	if err != nil {
		fmt.Println("cannot initialize", err)
		os.Exit(1)
	}
	return client
}

func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	return i
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}
