package main

import (
	"encoding/json"
	"fmt"
)

type J struct {
	ID   int
	Data []byte
}

type J2 struct {
	ID   int
	Data string
}

type J3 struct {
	ID   int
	Data [100]byte
}

func main() {
	data := `{"a":"b","c":"d","e":{"f":"g","h":{"i":"j"}}}`
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("marshal error: %v", err)
	}
	outJ := J{
		ID:   1,
		Data: b,
	}
	outJ2 := J2{
		ID:   1,
		Data: string(b),
	}
	outJ3 := J3{
		ID: 1,
	}
	copy(outJ3.Data[:], b)

	prettyPrint(outJ)
	prettyPrint(outJ2)
	prettyPrint(outJ3)
}

func prettyPrint(o interface{}) {
	outP, err := json.MarshalIndent(o, "", "	") // json.Marshal(o)
	if err != nil {
		fmt.Printf("out marshal error: %v", err)
	}
	fmt.Printf("%s\nlen:%d\n\n", string(outP), len(outP))
}
