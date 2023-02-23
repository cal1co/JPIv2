package dictconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cal1co/jpiv2/dboperations"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func AddEntries(IndexName string, client *opensearch.Client) {

	directory := "../jpiv2/dictdata/formatted/jmdict"

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFile, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}

			byteValue, _ := io.ReadAll(jsonFile)

			var entries []dboperations.Entry
			if err := json.Unmarshal(byteValue, &entries); err != nil {
				panic(err)
			}
			for i := 0; i < 300; i++ {
				fmt.Println(path, i)
				entry, err := json.Marshal(entries[i])
				if err != nil {
					log.Fatalf("Error marhaling JSON: %s", err.Error())
				}
				dboperations.Insert(
					opensearchapi.IndexRequest{
						Index: IndexName,
						Body:  strings.NewReader(string(entry)),
					},
					client)

			}

		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// func formatJMEntries() {
// var entryJson []dboperations.Entry
// for _, dictBank := range dict {
// 	jsonFile, err := os.Open(dictBank)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	byteValue, _ := io.ReadAll(jsonFile)

// 	var entries [][]interface{}
// 	if err := json.Unmarshal(byteValue, &entries); err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < len(entries); i++ {
// 		intId, err := strconv.Atoi(fmt.Sprint(entries[i][4]))
// 		if err != nil {
// 			panic(err)
// 		}
// 		def := entries[i][5].([]interface{})
// 		s := make([]string, len(def))
// 		for i, v := range def {
// 			s[i] = fmt.Sprint(v)
// 		}
// 		entry := dboperations.CreateEntry(fmt.Sprint(entries[i][0]), fmt.Sprint(entries[i][1]), fmt.Sprint(intId), s)
// 		entryJson = append(entryJson, entry)

// 		// dboperations.Insert(
// 		// 	opensearchapi.IndexRequest{
// 		// 		Index: IndexName,
// 		// 		Body:  entry,
// 		// 	},
// 		// 	client)
// 	}
// 	defer jsonFile.Close()
// }
// file, _ := json.MarshalIndent(entryJson, "", " ")
// os.WriteFile("jm_formatted_all.json", file, 0644)
// }
