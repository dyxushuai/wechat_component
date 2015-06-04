package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	a := &A{}
	fmt.Println(json.Unmarshal([]byte(`{}`), nil))
	fmt.Println(a.Name == "")
}

type A struct {
	Name string
}
