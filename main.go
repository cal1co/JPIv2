package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cal1co/jpiv2/dboperations"
	"github.com/cal1co/jpiv2/keygen"
	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
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
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/generatekey", func(w http.ResponseWriter, r *http.Request) {
		key, err := keygen.GenerateKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(key)
		fmt.Println(key)
		w.Header().Set("Content-Type", "application/json")

	})

	// Start the HTTP server.
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

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

func handleSearch(client *opensearch.Client, query string) *opensearchapi.Response {
	// Search
	search := opensearchapi.SearchRequest{
		Index: []string{IndexName},
		Body:  dboperations.CreateSearchQuery("50", query),
	}
	searchResponse := dboperations.Search(search, client)

	return searchResponse
}
