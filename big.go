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

// NewBig creates new Int by *big.Int.
// if `i` is nil then creates by new(big.Int)
func NewBig(i *big.Int) Int {
	if i == nil {
		i = new(big.Int)
	}
	return Int{Int: i}
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
		n := text[1 : len(text)-1]
		if bytes.HasPrefix(n, hprx) {
			r := string(n[2:])
			if i.Int, ok = new(big.Int).SetString(r, 16); !ok {
				return fmt.Errorf(`bigint: can't convert "0x%s" to *big.Int`, r)
			}
			return nil
		}

		r := string(n)
		if i.Int, ok = new(big.Int).SetString(r, 10); !ok {
			return fmt.Errorf(`bigint: can't convert "%s" to *big.Int`, r)
		}
		return nil
	}

	if bytes.Equal(text, null) {
		i.Int = new(big.Int)
		return nil
	}

	r := string(text)
	if i.Int, ok = new(big.Int).SetString(r, 10); !ok {
		return fmt.Errorf("bigint: can't convert %s to *big.Int", r)
	}
	return nil
}

// Scan implements the sql.Scanner interface.
// It converts decimal(N,0) or integer or NULL to *big.Int
// If the field is NULL,then creates by new(big.Int)
// 	var i Int
// 	_ = db.QueryRow("SELECT i FROM example WHERE id=1;").Scan(&i)
func (i *Int) Scan(val interface{}) error {
	if val == nil {
		i.Int = new(big.Int)
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
	case int32:
		i.Int = new(big.Int).SetInt64(int64(v))
		return nil
	case uint64:
		i.Int = new(big.Int).SetUint64(v)
		return nil
	case uint32:
		i.Int = new(big.Int).SetUint64(uint64(v))
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
//  _ = db.Exec("INSERT INTO example (i) VALUES (?);", bigint.New(100))
func (i Int) Value() (driver.Value, error) {
	if i.Int == nil {
		return "0", nil
	}
	return i.Int.String(), nil
}

// Copy creates new bigint.Int with deep copy
//  if `i.Int` is nil then creates by new(big.Int)
func (i Int) Copy() Int {
	if i.Int == nil {
		return Int{new(big.Int)}
	}
	return Int{new(big.Int).Set(i.Int)}
}
