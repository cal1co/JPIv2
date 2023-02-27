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
	"github.com/cal1co/jpiv2/keygen"
	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	"golang.org/x/time/rate"
)

const IndexName = "go-test-index1"

type Entry struct {
	Word      string
	Alternate string
	Freq      string
	Def       []string
}
type Hit struct {
	// Id     string  `json:"_id`
	Score  float64 `json:"_score"`
	Index  string  `json:"_index"`
	Source Entry   `json:"_source"`
}
type Hits struct {
	Total struct {
		Value int64
	}
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}
type SearchResult struct {
	Took     int
	TimedOut bool `json:"timed_out"`
	Hits     Hits
}

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

var limiter = NewIPRateLimiter(1, 5)

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
		// w.Header().Set("Content-Type", "application/json")
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if len(query) > 0 {
			searchRes := handleSearch(client, query)

			var searchResult SearchResult
			if err := json.NewDecoder(searchRes.Body).Decode(&searchResult); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(os.Stdout).Encode(searchResult)
			if err := json.NewEncoder(w).Encode(searchResult.Hits); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	})

	mux.HandleFunc("/generatekey", generateKeyHandler)

	// Start the HTTP server.
	fmt.Println("Server listening on port 8888")

	if err := http.ListenAndServe(":8888", limitMiddleware(mux)); err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	// Delete the previously created index.
	// deleteIndex := opensearchapi.IndicesDeleteRequest{
	// 	Index: []string{IndexName},
	// }

	// deleteIndexResponse, err := deleteIndex.Do(context.Background(), client)
	// if err != nil {
	// 	fmt.Println("failed to delete index ", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Deleting the index")
	// fmt.Println(deleteIndexResponse)
	// defer deleteIndexResponse.Body.Close()
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

func handleSearch(client *opensearch.Client, query string) *opensearchapi.Response {
	// Search
	search := opensearchapi.SearchRequest{
		Index: []string{IndexName},
		Body:  dboperations.CreateSearchQuery("50", query),
	}
	searchResponse := dboperations.Search(search, client)

	return searchResponse
}

func generateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key, err := keygen.GenerateKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(key)
	fmt.Println(key)
	w.Header().Set("Content-Type", "application/json")
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
