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

	// null
	var a Object
	_ = json.Unmarshal([]byte(`{"Field": null}`), &a)
	fmt.Println(a.Field)

	// unsigned integer
	var b Object
	_ = json.Unmarshal([]byte(`{"Field": 1}`), &b)
	fmt.Println(b.Field)

	// string unsigned integer
	var c Object
	_ = json.Unmarshal([]byte(`{"Field": "2"}`), &c)
	fmt.Println(c.Field)

	// hex string
	var d Object
	_ = json.Unmarshal([]byte(`{"Field": "0x3"}`), &d)
	fmt.Println(d.Field)

	// not supports decoding emtpy string or empty hex
	// eg. `{"Field": ""}` or `{"Field": "0x"}`

	// Output:
	// 0
	// 1
	// 2
	// 3
}

func ExampleNew() {
	var i = New(100)
	i.Add(i.Int, big.NewInt(100))
	fmt.Println(i)

	// Output:
	// 200
}

func ExampleNewBig() {
	var a = NewBig(nil)
	fmt.Println(a)

	var b = NewBig(big.NewInt(1))
	fmt.Println(b)

	// Output:
	// 0
	// 1
}
