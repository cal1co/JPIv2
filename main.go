package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cal1co/jpiv2/dboperations"
	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

const IndexName = "go-test-index1"

func main() {
	// Initialize the client with SSL/TLS enabled.
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

	// Print OpenSearch version information on console.
	fmt.Println(client.Info())

	// Define index settings.
	settings := strings.NewReader(`{
     'settings': {
       'index': {
            'number_of_shards': 1,
            'number_of_replicas': 2
            }
          }
     }`)

	// Create an index with non-default settings.
	res := opensearchapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  settings,
	}
	fmt.Println("Creating index")
	fmt.Println(res)

	// Insert
	entry_document := dboperations.CreateEntry("Moneyball", "Bennett Miller", "2011")
	req := opensearchapi.IndexRequest{
		Index: IndexName,
		Body:  entry_document,
	}
	insertResponse := dboperations.Insert(req, client)
	defer insertResponse.Body.Close()

	// Search
	query := dboperations.CreateSearchQuery("5", "Moneyball")
	search := opensearchapi.SearchRequest{
		Index: []string{IndexName},
		Body:  query,
	}
	searchResponse := dboperations.Search(search, client)
	fmt.Println("SEARCH RESULTS: ", searchResponse)
	defer searchResponse.Body.Close()

	// Delete the previously created index.
	deleteIndex := opensearchapi.IndicesDeleteRequest{
		Index: []string{IndexName},
	}

	deleteIndexResponse, err := deleteIndex.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to delete index ", err)
		os.Exit(1)
	}
	fmt.Println("Deleting the index")
	fmt.Println(deleteIndexResponse)
	defer deleteIndexResponse.Body.Close()
}
