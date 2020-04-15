package bigint

import (
	"encoding/json"
	"reflect"
	"testing"
)

func BenchmarkInt_UnmarshalJSON(b *testing.B) {
	type Object struct {
		String  Int
		Integer Int
		Null    Int
		Hex     Int
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
