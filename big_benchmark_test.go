package bigint

import (
	"encoding/json"
	"reflect"
	"testing"
)

/*
BenchmarkInt_UnmarshalJSON-12    	  511303	      2218 ns/op	    1064 B/op	      21 allocs/op
BenchmarkInt_UnmarshalJSON-12    	  538291	      2205 ns/op	    1064 B/op	      21 allocs/op
BenchmarkInt_UnmarshalJSON-12    	  753636	      1589 ns/op	     488 B/op	      15 allocs/op
BenchmarkInt_UnmarshalJSON-12    	  764230	      1581 ns/op	     480 B/op	      14 allocs/op
BenchmarkInt_UnmarshalJSON-12    	  764990	      1553 ns/op	     480 B/op	      14 allocs/op
*/

func BenchmarkInt_UnmarshalJSON(b *testing.B) {
	type Object struct {
		Null    Int
		Hex     Int
		String  Int
		Integer Int
	}

	var raw = []byte(`{"Null":null,"Hex":"0x1","String":"2","Integer":"3"}`)
	want := &Object{Null: New(0), Hex: New(1), String: New(2), Integer: New(3)}

	var O Object
	if err := json.Unmarshal(raw, &O); err != nil {
		b.Errorf("decode error %s", err)
	}

	if !reflect.DeepEqual(&O, want) {
		b.Fatalf("Not equal: want %+v but got %+v", want, &O)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var O Object
		_ = json.Unmarshal(raw, &O)
	}
}
