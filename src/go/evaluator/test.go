package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	options := make(map[string]interface{})
	a := `{"age": 12, "name": "sevenpan"}`
	_ = json.Unmarshal([]byte(a), &options)
	fmt.Println(options["age"])
	age := options["age"]
	ageNumber, ok := age.(json.Number)
	if ok {
		fmt.Println("is number")
		if _, err := ageNumber.Int64(); err == nil {
			fmt.Println("is int64")
		}

	}
}
