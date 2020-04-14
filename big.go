package bigint

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

var quote = []byte(`"`)
var null = []byte(`null`)

const hexprx = "0x"

type Int struct {
	*big.Int
}

// New creates Int by uint64
func New(i uint64) Int {
	return Int{Int: new(big.Int).SetUint64(i)}
}

// NewBig creates Int by *big.Int
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
		// shouldn't use null variable on above
		return []byte("null"), nil
	}
	return json.Marshal(i.Int.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It converts integer or string integer or hex string to *big.Int.
func (i *Int) UnmarshalJSON(text []byte) error {
	if bytes.Equal(text, null) {
		i.Int = new(big.Int)
		return nil
	}

	if bytes.HasPrefix(text, quote) && bytes.HasSuffix(text, quote) {
		var s string
		_ = json.Unmarshal(text, &s) // always no error

		var ok bool
		if strings.HasPrefix(s, hexprx) {
			if i.Int, ok = new(big.Int).SetString(s[2:], 16); !ok {
				return fmt.Errorf("bigint: can't convert hex %s to *big.Int", string(text))
			}
			return nil
		}

		if i.Int, ok = new(big.Int).SetString(s, 10); !ok {
			return fmt.Errorf("bigint: can't convert %s to *big.Int", string(text))
		}
		return nil
	}

	var x int64
	if err := json.Unmarshal(text, &x); err != nil {
		return err
	}
	i.Int = new(big.Int).SetInt64(x)
	return nil
}

// Scan implements the sql.Scanner interface.
// It converts decimal(N,0) or integer or NULL to *big.Int
// If the field is NULL,bigint.Int will create by new(big.Int)
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
	case uint64:
		i.Int = new(big.Int).SetUint64(v)
		return nil
	case int32:
		i.Int = new(big.Int).SetInt64(int64(v))
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
//  var i = Int{big.NewInt(100)}
//  _ = db.Exec("INSERT INTO example (i) VALUES (?);", i)
func (i Int) Value() (driver.Value, error) {
	if i.Int == nil {
		return "0", nil
	}
	return i.Int.String(), nil
}

// Copy create new *big.Int with deep copy
func (i Int) Copy() Int {
	if i.Int == nil {
		return Int{new(big.Int)}
	}
	return Int{new(big.Int).Set(i.Int)}
}
