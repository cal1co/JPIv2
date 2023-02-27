package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cal1co/jpiv2/dboperations"
	"github.com/cal1co/jpiv2/keygen"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

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

func GenerateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key, err := keygen.GenerateKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(key)
	fmt.Println(key)
	w.Header().Set("Content-Type", "application/json")
}

func SearchHandler(w http.ResponseWriter, r *http.Request, client *opensearch.Client, IndexName string) {
	query := r.URL.Query().Get("query")
	if len(query) > 0 {
		searchRes := handleSearch(client, query, IndexName)

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
}

func handleSearch(client *opensearch.Client, query string, IndexName string) *opensearchapi.Response {
	// Search
	search := opensearchapi.SearchRequest{
		Index: []string{IndexName},
		Body:  dboperations.CreateSearchQuery("50", query),
	}
	searchResponse := dboperations.Search(search, client)

	return searchResponse
}
