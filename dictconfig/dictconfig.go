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

type Entry struct {
	Word      string
	Alternate string
	Freq      string
	Def       []string
	Pitch     string
}

type EntryGroup struct {
	Entries []Entry
	Pitch   string
}

func AddEntries(IndexName string, client *opensearch.Client) {

	directory := "../jpiv2/dictdata/formatted/jm_formatted_all_pitch.json"

	jsonFile, err := os.Open(directory)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var entries []Entry
	if err := json.Unmarshal(byteValue, &entries); err != nil {
		panic(err)
	}
	for i := 20000; i < 30000; i++ {
		fmt.Println(i)
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

	// err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !info.IsDir() && filepath.Ext(path) == ".json" {
	// 		jsonFile, err := os.Open(path)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}

	// 		byteValue, _ := io.ReadAll(jsonFile)

	// 		var entries []dboperations.Entry
	// 		if err := json.Unmarshal(byteValue, &entries); err != nil {
	// 			panic(err)
	// 		}
	// 		for i := 0; i < len(entries); i++ {
	// 			fmt.Println(path, i)
	// 			entry, err := json.Marshal(entries[i])
	// 			if err != nil {
	// 				log.Fatalf("Error marhaling JSON: %s", err.Error())
	// 			}
	// 			dboperations.Insert(
	// 				opensearchapi.IndexRequest{
	// 					Index: IndexName,
	// 					Body:  strings.NewReader(string(entry)),
	// 				},
	// 				client)

	// 		}

	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func AddPitchToJM() {
	jmDict := make(map[string]EntryGroup)
	jm_directory := "../jpiv2/dictdata/formatted/jmdict"
	err := filepath.Walk(jm_directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFile, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}

			byteValue, _ := io.ReadAll(jsonFile)
			var entries []Entry
			if err := json.Unmarshal(byteValue, &entries); err != nil {
				panic(err)
			}
			for i := 0; i < len(entries); i++ {
				if _, ok := jmDict[entries[i].Word]; ok {
					entryStruct := jmDict[entries[i].Word]
					entryStruct.Entries = append(jmDict[entries[i].Word].Entries, entries[i])
					jmDict[entries[i].Word] = entryStruct
				} else {
					var entry EntryGroup
					entry.Entries = append(entry.Entries, entries[i])
					jmDict[entries[i].Word] = entry
				}
			}
			defer jsonFile.Close()
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	directory := "../jpiv2/dictdata/NHK_accent/"
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFile, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}

			byteValue, _ := io.ReadAll(jsonFile)
			var entries [][]interface{}
			if err := json.Unmarshal(byteValue, &entries); err != nil {
				panic(err)
			}
			for i := 0; i < len(entries); i++ {
				word := fmt.Sprintf(entries[i][0].(string))
				if _, ok := jmDict[word]; ok {
					if wordStruct, ok := jmDict[word]; ok {
						pitchInfoInterface := entries[i][5].([]interface{})
						var pitchInfoString []string
						for _, val := range pitchInfoInterface {
							str, ok := val.(string)
							if ok {
								pitchInfoString = append(pitchInfoString, str)
							}
						}
						wordStruct.Pitch = pitchInfoString[0]
						var entryList []Entry
						for _, entry := range wordStruct.Entries {
							entry.Pitch = pitchInfoString[0]
							entryList = append(entryList, entry)
						}
						wordStruct.Entries = entryList
						jmDict[word] = wordStruct

					}
				}
			}
			defer jsonFile.Close()
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	var entryJson []Entry
	for _, val := range jmDict {
		entryJson = append(entryJson, val.Entries...)
	}
	file, _ := json.MarshalIndent(entryJson, "", " ")
	os.WriteFile("jm_formatted_all_pitch.json", file, 0644)
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
