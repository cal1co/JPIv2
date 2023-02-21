package dboperations

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func CreateEntry(title string, direction string, year string) *strings.Reader {
	reader := `{
        "title": "%s",
        "director": "%s",
        "year": "%s"
    }`
	formattedStr := fmt.Sprintf(reader, title, direction, year)
	return strings.NewReader(formattedStr)
}

func Insert(req opensearchapi.IndexRequest, client *opensearch.Client) *opensearchapi.Response {
	response, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println("OOPS ERROR IN INSERT")
	}
	return response
}

func CreateSearchQuery(size string, query string) *strings.Reader {
	reader := `{
		"size": %s,
		"query": {
			"multi_match": {
			"query": "%s",
			"fields": ["title", "director"]
			}
	   }
	 }`
	formattedStr := fmt.Sprintf(reader, size, query)
	// fmt.Println("FORMATTED SEARCH QUERY HERE: ", formattedStr)
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

func DeleteSpecific() {
	// Delete the document.
	// docId := "1"
	// delete := opensearchapi.DeleteRequest{
	// 	Index:      IndexName,
	// 	DocumentID: docId,
	// }

	// deleteResponse, err := delete.Do(context.Background(), client)
	// if err != nil {
	// 	fmt.Println("failed to delete document ", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Deleting a document")
	// fmt.Println(deleteResponse)
	// defer deleteResponse.Body.Close()

}

func DropAll() {

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
