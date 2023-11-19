package main

import (
	"encoding/json"
	"fmt"
)

type test struct {
	Name string
	Age  int
}

func main() {
	var test_slice []test
	for i := 0; i < 3; i++ {
		test1 := test{
			Name: "123",
			Age:  1,
		}
		test_slice = append(test_slice, test1)
	}
	fmt.Println(test_slice)
	data, _ := json.Marshal(test_slice)

	var test_slice_unmarshal []test
	_ = json.Unmarshal(data, &test_slice_unmarshal)
	fmt.Println(test_slice_unmarshal)
}
