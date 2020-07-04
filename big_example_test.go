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
	fmt.Println(a.Field, a.Field.Int == nil)

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

	var e Object
	_ = json.Unmarshal([]byte(`{"Field": "0x"}`), &e)
	fmt.Println(e.Field, e.Field.Int == nil)

	var f Object
	_ = json.Unmarshal([]byte(`{"Field": ""}`), &f)
	fmt.Println(f.Field, f.Field.Int == nil)

	// not supports decoding emtpy string or empty hex
	// eg. `{"Field": ""}` or `{"Field": "0x"}`

	// Output:
	// <nil> true
	// 1
	// 2
	// 3
	// <nil> true
	// <nil> true
}

func ExampleNew() {
	var i = New(100)
	i.Add(i.Int, big.NewInt(100))
	j := i.Copy()

	i.Add(i.Int, big.NewInt(100))
	fmt.Println(i, j)

	// Output:
	// 300 200
}
