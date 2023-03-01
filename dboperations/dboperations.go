package dboperations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func CreateIndex(IndexName string) opensearchapi.IndicesCreateRequest {
	settings := strings.NewReader(`{
		'settings': {
		  'index': {
			   'number_of_shards': 5,
			   'number_of_replicas': 2
			   }
			 }
		}`)
	res := opensearchapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  settings,
	}
	return res
}

type Entry struct {
	Word      string
	Alternate string
	Freq      string
	Def       []string
}

func CreateEntry(word string, alternate string, freq string, def []string) *strings.Reader {

	entry := Entry{
		Word:      word,
		Alternate: alternate,
		Freq:      freq,
		Def:       def,
	}

	jsonStr, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("Error marhaling JSON: %s", err.Error())
	}
	return strings.NewReader(string(jsonStr))
}

func Insert(req opensearchapi.IndexRequest, client *opensearch.Client) *opensearchapi.Response {
	response, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println("OOPS ERROR IN INSERT", err)
		return nil
	}
	return response
}

func CreateSearchQuery(size string, query string) *strings.Reader {
	reader := `{
		"size": %s,
		"query": {
			"multi_match": {
			"query": "%s",
			"fields": ["Def", "Word", "Alternate"]
			}
	   }
	 }`
	formattedStr := fmt.Sprintf(reader, size, query)
	return strings.NewReader(formattedStr)
}

func Search(search opensearchapi.SearchRequest, client *opensearch.Client) *opensearchapi.Response {
	res, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to search document ", err)
		os.Exit(1)
	}
	return res
}

func DeleteSpecific(IndexName string, docId string, client *opensearch.Client) {
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
}

func Cleanup(IndexName string, client *opensearch.Client) {
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

func PerformBulk() {
	// Perform bulk operations.
	// 	blk, err := client.Bulk(
	// 		strings.NewReader(`
	// 	{ "index" : { "_index" : "go-test-index1" } }
	// 	{ "title" : "Interstellar", "director" : "Christopher Nolan", "year" : "2014"}
	// 	{ "create" : { "_index" : "go-test-index1" } }
	// 	{ "title" : "Star Trek Beyond", "director" : "Justin Lin", "year" : "2015"}
	// 	{ "update" : {"_id" : "3", "_index" : "go-test-index1" } }
	// 	{ "doc" : {"year" : "2016"} }
	// `),
	// 	)

	// if err != nil {
	// 	fmt.Println("failed to perform bulk operations", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Performing bulk operations")
	// fmt.Println(blk)
}
