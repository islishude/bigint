package bigint

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func ExampleInt_UnmarshalJSON() {
	type Object struct {
		Field Int
	}

	var a Object
	_ = json.Unmarshal([]byte(`{"Field": null}`), &a)
	fmt.Println(a.Field)

	var b Object
	_ = json.Unmarshal([]byte(`{"Field": 1}`), &b)
	fmt.Println(b.Field)

	var c Object
	_ = json.Unmarshal([]byte(`{"Field": "2"}`), &c)
	fmt.Println(c.Field)

	// Output:
	// 0
	// 1
	// 2
}

func ExampleNew() {
	var i = New(big.NewInt(100))
	i.Add(i.Int, big.NewInt(100))
	fmt.Println(i)

	// Output:
	// 200
}
