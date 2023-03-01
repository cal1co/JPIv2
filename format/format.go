package format

import "fmt"

type Entry struct {
	Word      string
	Alternate string
	Freq      string
	Def       []string
	Pitch     string
}

func AddPitchToJM() {
	fmt.Println("ADD PITCH CALLED")
}
