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

const IndexName = "jpiv2-dict-store"

var limiter = NewIPRateLimiter(1, 5)

type IPRateLimiter struct {
	ipMap map[string]*rate.Limiter
	mu    *sync.RWMutex
	r     rate.Limit
	b     int
}

func main() {

	client := initClient()
	fmt.Println(client.Info())

	res := dboperations.CreateIndex(IndexName)
	fmt.Println("Creating index")
	fmt.Println(res)

	// dictconfig.AddEntries(IndexName, client)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("OK")
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.SearchHandler(w, r, client, IndexName)
	})

	mux.HandleFunc("/generatekey", handlers.GenerateKeyHandler)

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
		ipMap: make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		r:     r,
		b:     b,
	}

	return i
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ipMap[ip] = limiter

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ipMap[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}
