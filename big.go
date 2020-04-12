package bigint

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
)

var quote = []byte(`"`)
var null = []byte(`null`)

type Int struct {
	*big.Int
}

// New copies *big.Int to Int
// If i is nil and then create new *big.Int
func New(i *big.Int) Int {
	if i == nil {
		return Int{Int: new(big.Int)}
	}
	return Int{Int: new(big.Int).Set(i)}
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
// It supports unsigned integer(uint64) or string integer to *big.Int.
func (i *Int) UnmarshalJSON(text []byte) error {
	if bytes.Equal(text, null) {
		i.Int = new(big.Int)
		return nil
	}

	if bytes.HasPrefix(text, quote) && bytes.HasSuffix(text, quote) {
		var s string
		_ = json.Unmarshal(text, &s) // always no error

		var ok bool
		if i.Int, ok = new(big.Int).SetString(s, 10); !ok {
			return fmt.Errorf("Can't convert %s to *big.Int", string(text))
		}
		return nil
	}

	var x uint64
	if err := json.Unmarshal(text, &x); err != nil {
		return err
	}
	i.Int = new(big.Int).SetUint64(x)
	return nil
}

// Scan implements the sql.Scanner interface.
// It converts decimal(N,0) to *big.Int or NULL to new(big.Int)
// Example:
// 	var i Int
// 	_ = db.QueryRow("SELECT i FROM example WHERE id=1;").Scan(&i)
func (i *Int) Scan(val interface{}) error {
	if val == nil {
		i.Int = new(big.Int)
		return nil
	}

	var data string
	switch i := val.(type) {
	case []byte:
		data = string(i)
	case string:
		data = i
	default:
		return fmt.Errorf("Can't Scan to *big.Int by %s", reflect.TypeOf(val).Kind())
	}

	var ok bool
	i.Int, ok = new(big.Int).SetString(data, 10)
	if !ok {
		return fmt.Errorf("Can't convert %s to *big.Int", data)
	}
	return nil
}

// Value implements the driver.Valuer interface.
// Example:
//  var i = Int{big.NewInt(100)}
//	_ = db.Exec("INSERT INTO example (i) VALUES (?);", i)
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
