package dictconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/cal1co/jpiv2/dboperations"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type Entry struct {
	Word      string
	Alternate string
	Freq      string
	Def       []string
	Dist      int
}

func AddEntries(IndexName string, client *opensearch.Client) {

	dict := []string{
		"dictdata/jmdict_english/term_bank_1.json",
		// "dictdata/jmdict_english/term_bank_2.json",
		// "dictdata/jmdict_english/term_bank_3.json",
		// "dictdata/jmdict_english/term_bank_4.json",
		// "dictdata/jmdict_english/term_bank_5.json",
		// "dictdata/jmdict_english/term_bank_6.json",
		// "dictdata/jmdict_english/term_bank_7.json",
		// "dictdata/jmdict_english/term_bank_8.json",
		// "dictdata/jmdict_english/term_bank_9.json",
		// "dictdata/jmdict_english/term_bank_10.json",
		// "dictdata/jmdict_english/term_bank_11.json",
		// "dictdata/jmdict_english/term_bank_12.json",
		// "dictdata/jmdict_english/term_bank_13.json",
		// "dictdata/jmdict_english/term_bank_14.json",
		// "dictdata/jmdict_english/term_bank_15.json",
		// "dictdata/jmdict_english/term_bank_16.json",
		// "dictdata/jmdict_english/term_bank_17.json",
		// "dictdata/jmdict_english/term_bank_18.json",
		// "dictdata/jmdict_english/term_bank_19.json",
		// "dictdata/jmdict_english/term_bank_20.json",
		// "dictdata/jmdict_english/term_bank_21.json",
		// "dictdata/jmdict_english/term_bank_22.json",
		// "dictdata/jmdict_english/term_bank_23.json",
		// "dictdata/jmdict_english/term_bank_24.json",
		// "dictdata/jmdict_english/term_bank_25.json",
		// "dictdata/jmdict_english/term_bank_26.json",
		// "dictdata/jmdict_english/term_bank_27.json",
		// "dictdata/jmdict_english/term_bank_28.json",
		// "dictdata/jmdict_english/term_bank_29.json",
	}

	for _, dictBank := range dict {
		jsonFile, err := os.Open(dictBank)
		if err != nil {
			fmt.Println(err)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var entries [][]interface{}
		if err := json.Unmarshal(byteValue, &entries); err != nil {
			panic(err)
		}

		for i := 0; i < len(entries); i++ {
			fmt.Println(i)
			// for i := 0; i < 500; i++ {
			intId, err := strconv.Atoi(fmt.Sprint(entries[i][4]))
			if err != nil {
				panic(err)
			}
			def := entries[i][5].([]interface{})
			s := make([]string, len(def))
			for i, v := range def {
				s[i] = fmt.Sprint(v)
			}
			entry := dboperations.CreateEntry(fmt.Sprint(entries[i][0]), fmt.Sprint(entries[i][1]), fmt.Sprint(intId), s)

			dboperations.Insert(
				opensearchapi.IndexRequest{
					Index: IndexName,
					Body:  entry,
				},
				client)
		}
	}
}
