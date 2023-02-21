package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	// Add a document to the index.
	document := strings.NewReader(`{
        "title": "Moneyball",
        "director": "Bennett Miller",
        "year": "2011"
    }`)

	docId := "1"
	req := opensearchapi.IndexRequest{
		Index:      IndexName,
		DocumentID: docId,
		Body:       document,
	}
	insertResponse, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to insert document ", err)
		os.Exit(1)
	}
	fmt.Println("Inserting a document")
	fmt.Println(insertResponse)
	defer insertResponse.Body.Close()

	// Perform bulk operations.
	blk, err := client.Bulk(
		strings.NewReader(`
    { "index" : { "_index" : "go-test-index1", "_id" : "2" } }
    { "title" : "Interstellar", "director" : "Christopher Nolan", "year" : "2014"}
    { "create" : { "_index" : "go-test-index1", "_id" : "3" } }
    { "title" : "Star Trek Beyond", "director" : "Justin Lin", "year" : "2015"}
    { "update" : {"_id" : "3", "_index" : "go-test-index1" } }
    { "doc" : {"year" : "2016"} }
`),
	)

	if err != nil {
		fmt.Println("failed to perform bulk operations", err)
		os.Exit(1)
	}
	fmt.Println("Performing bulk operations")
	fmt.Println(blk)

	// Search for the document.
	content := strings.NewReader(`{
       "size": 5,
       "query": {
           "multi_match": {
           "query": "Justin Lin",
           "fields": ["title", "director"]
           }
      }
    }`)

	search := opensearchapi.SearchRequest{
		Index: []string{IndexName},
		Body:  content,
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to search document ", err)
		os.Exit(1)
	}
	fmt.Println("Searching for a document")
	fmt.Println(searchResponse)
	defer searchResponse.Body.Close()

	// Delete the document.
	delete := opensearchapi.DeleteRequest{
		Index:      IndexName,
		DocumentID: docId,
	}

	deleteResponse, err := delete.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to delete document ", err)
		os.Exit(1)
	}
	fmt.Println("Deleting a document")
	fmt.Println(deleteResponse)
	defer deleteResponse.Body.Close()

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
