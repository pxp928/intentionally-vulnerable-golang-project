package main

import (
	"github.com/mailru/easyjson"
)

//easyjson:json
type SomeStruct struct {
	Field1 string
	Field2 string
}

func main() {
	someStruct := SomeStruct{Field1: "val1", Field2: "val2"}
	_, _ = easyjson.Marshal(someStruct)
}
