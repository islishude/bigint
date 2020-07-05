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

	// signed integer
	_ = json.Unmarshal([]byte(`{"Field": -1}`), &b)
	fmt.Println(b.Field)

	// string unsigned integer
	var c Object
	_ = json.Unmarshal([]byte(`{"Field": "2"}`), &c)
	fmt.Println(c.Field)

	_ = json.Unmarshal([]byte(`{"Field": "-2"}`), &c)
	fmt.Println(c.Field)

	// hex string
	var d Object
	_ = json.Unmarshal([]byte(`{"Field": "0x3"}`), &d)
	fmt.Println(d.Field)

	// empty hex string
	var e Object
	_ = json.Unmarshal([]byte(`{"Field": "0x"}`), &e)
	fmt.Println(e.Field)

	// empty string
	var f Object
	_ = json.Unmarshal([]byte(`{"Field": ""}`), &f)
	fmt.Println(f.Field)

	// Output:
	// <nil>
	// 1
	// -1
	// 2
	// -2
	// 3
	// <nil>
	// <nil>
}

func ExampleInt_Copy() {
	i := New(100)
	j := i.Copy()

	i.Add(i.Int, big.NewInt(100))

	fmt.Println(i, j)
	// Output:
	// 200 100
}

func ExampleInt_Safer() {
	var a Int
	_ = json.Unmarshal([]byte(`{"Field": null}`), &a)

	fmt.Println(a.IsNil())
	func() {
		defer func() {
			_ = recover()
			fmt.Println("panic!")
		}()
		a.Add(big.NewInt(100), big.NewInt(100))
	}()

	fmt.Println(a.Safer().IsNil())
	fmt.Println(a.IsNil()) // It's already safe

	a.Add(big.NewInt(100), big.NewInt(100))
	fmt.Println(a)

	// Output:
	// true
	// panic!
	// false
	// false
	// 200
}
