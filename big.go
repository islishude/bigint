package bigint

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
)

var (
	quote = []byte{34}                 // `"`
	hprx  = []byte{48, 120}            // `0x`
	null  = []byte{110, 117, 108, 108} // `null`
)

type Int struct {
	*big.Int
}

// New creates Int by int64
func New(i int64) Int {
	return Int{Int: big.NewInt(i)}
}

// New creates Int by int64
func NewUint(i uint64) Int {
	return Int{Int: new(big.Int).SetUint64(i)}
}

// MarshalJSON implements the json.Marshaler interface.
// It converts *big.Int to string integer
func (i Int) MarshalJSON() ([]byte, error) {
	if i.Int == nil {
		return append([]byte(nil), null...), nil
	}
	return json.Marshal(i.Int.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It converts integer or string integer or hex string to *big.Int.
func (i *Int) UnmarshalJSON(text []byte) error {
	var ok bool
	if bytes.HasPrefix(text, quote) {
		text = text[1 : len(text)-1]
	}

	if bytes.HasPrefix(text, hprx) {
		if r := string(text[2:]); r != "" {
			if i.Int, ok = new(big.Int).SetString(r, 16); !ok {
				return fmt.Errorf(`bigint: can't convert "0x%s" to *big.Int`, r)
			}
		}
		return nil
	}

	if bytes.Equal(text, null) {
		return nil
	}

	if r := string(text); r != "" {
		if i.Int, ok = new(big.Int).SetString(r, 10); !ok {
			return fmt.Errorf(`bigint: can't convert "%s" to *big.Int`, r)
		}
	}
	return nil
}

// Scan implements the sql.Scanner interface.
// It converts decimal(N,0) or integer or NULL to *big.Int
// 	var i Int
// 	_ = db.QueryRow("SELECT i FROM example WHERE id=1;").Scan(&i)
func (i *Int) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	var data string
	switch v := val.(type) {
	case []byte:
		data = string(v)
	case string:
		data = v
	case int64:
		i.Int = new(big.Int).SetInt64(v)
		return nil
	default:
		return fmt.Errorf("bigint: can't convert %s type to *big.Int", reflect.TypeOf(val).Kind())
	}

	var ok bool
	i.Int, ok = new(big.Int).SetString(data, 10)
	if !ok {
		return fmt.Errorf("bigint can't convert %s to *big.Int", data)
	}
	return nil
}

// Value implements the driver.Valuer interface.
//  var i bigint.Int
//  // only for nullable field
//  _ = db.Exec("INSERT INTO example (i) VALUES (?);", i)
//  i = bigint.New(1024)
//  _ = db.Exec("INSERT INTO example (i) VALUES (?);", i)
func (i Int) Value() (driver.Value, error) {
	if i.Int == nil {
		return nil, nil
	}
	return i.Int.String(), nil
}

// Copy creates new bigint.Int with deeply copy
func (i Int) Copy() Int {
	if i.Int == nil {
		return Int{}
	}
	return Int{new(big.Int).Set(i.Int)}
}

// IsNil returns is or not nil
func (i Int) IsNil() bool {
	return i.Int == nil
}

// Safer converts nil value to new(big.Int)
func (i *Int) Safer() *Int {
	if i.Int == nil {
		i.Int = new(big.Int)
	}
	return i
}
