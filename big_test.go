package bigint

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
)

func TestInt_MarshalJSON(t *testing.T) {
	type fields struct {
		Int *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "nil",
			fields:  fields{},
			want:    []byte(`{"Field":null}`),
			wantErr: false,
		},
		{
			name:    "not nil",
			fields:  fields{big.NewInt(1)},
			want:    []byte(`{"Field":"1"}`),
			wantErr: false,
		},
	}

	type Object struct {
		Field Int
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(Object{Field: Int{tt.fields.Int}})
			if (err != nil) != tt.wantErr {
				t.Errorf("Int.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Int.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestInt_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Int *big.Int
	}
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *big.Int
	}{
		{
			name:    "invalid json",
			args:    args{[]byte(`{"Field": "1}`)},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid json 2",
			args:    args{[]byte(`{"Field": 1"}`)},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "null",
			args:    args{[]byte(`{"Field": null}`)},
			wantErr: false,
			want:    new(big.Int),
		},
		{
			name:    "string integer",
			args:    args{[]byte(`{"Field": "1024"}`)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "negative string integer",
			args:    args{[]byte(`{"Field": "-1024"}`)},
			wantErr: false,
			want:    big.NewInt(-1024),
		},
		{
			name:    "string hex",
			args:    args{[]byte(`{"Field": "0x400"}`)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "string hex 2",
			args:    args{[]byte(`{"Field": "0x03"}`)},
			wantErr: false,
			want:    big.NewInt(3),
		},
		{
			name:    "invalid hex string",
			args:    args{[]byte(`{"Field": "0xxyz"}`)},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "integer",
			args:    args{[]byte(`{"Field": 1024}`)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "negative integer",
			args:    args{[]byte(`{"Field": -1024}`)},
			wantErr: false,
			want:    big.NewInt(-1024),
		},
		{
			name:    "float64 string",
			args:    args{[]byte(`{"Field": "10.24"}`)},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "float64",
			args:    args{[]byte(`{"Field": 10.24}`)},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "not a number",
			args:    args{[]byte(`{"Field": "abc"}`)},
			wantErr: true,
			want:    nil,
		},
	}

	type Object struct {
		Field Int
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Object
			if err := json.Unmarshal(tt.args.text, &o); (err != nil) != tt.wantErr {
				t.Errorf("Int.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(o.Field.Int, tt.want) {
				t.Errorf("Int.UnmarshalJSON() want %s, got %s", tt.want.String(), o.Field.String())
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		i uint64
	}
	tests := []struct {
		name string
		args args
		want Int
	}{
		{
			name: "nil",
			args: args{1024},
			want: Int{big.NewInt(1024)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt_Copy(t *testing.T) {
	type fields struct {
		Int *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   Int
	}{
		{
			name:   "nil",
			fields: fields{},
			want:   Int{new(big.Int)},
		},
		{
			name:   "not nil",
			fields: fields{big.NewInt(1024)},
			want:   Int{Int: big.NewInt(1024)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Int{
				Int: tt.fields.Int,
			}
			if got := i.Copy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int.Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt_Value(t *testing.T) {
	type fields struct {
		Int *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "nil",
			fields:  fields{},
			want:    "0",
			wantErr: false,
		},
		{
			name:    "not nil",
			fields:  fields{big.NewInt(1024)},
			want:    "1024",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Int{
				Int: tt.fields.Int,
			}
			got, err := i.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt_Scan(t *testing.T) {
	type fields struct {
		Int *big.Int
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *big.Int
	}{
		{
			name:    "valid string",
			fields:  fields{nil},
			args:    args{"1024"},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "valid string with init value",
			fields:  fields{big.NewInt(1)},
			args:    args{"1024"},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "invalid decimal string",
			fields:  fields{},
			args:    args{"abc"},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "valid bytes",
			fields:  fields{},
			args:    args{[]byte("1024")},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "decimal(10,2)",
			fields:  fields{},
			args:    args{[]byte("10.24")},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "double",
			fields:  fields{},
			args:    args{10.24},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "valid int type",
			fields:  fields{},
			args:    args{100},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "null type",
			fields:  fields{},
			args:    args{nil},
			wantErr: false,
			want:    new(big.Int),
		},
		{
			name:    "int32",
			fields:  fields{nil},
			args:    args{int32(1024)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "uint32",
			fields:  fields{nil},
			args:    args{uint32(1024)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "int64",
			fields:  fields{nil},
			args:    args{int64(1024)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
		{
			name:    "uint64",
			fields:  fields{nil},
			args:    args{uint64(1024)},
			wantErr: false,
			want:    big.NewInt(1024),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int{
				Int: tt.fields.Int,
			}
			if err := i.Scan(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Int.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(i.Int, tt.want) {
				t.Errorf("Int.Scan() want %s, got %s", tt.want, i.Int)
			}
		})
	}
}
