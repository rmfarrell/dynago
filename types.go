package dynago

import (
	"encoding/json"
	"strconv"
)

type StringSet []string

type NumberSet []string

type BinarySet [][]byte

type List []interface{}

type Number string

func (n Number) IntVal() (int, error) {
	return strconv.Atoi(string(n))
}

func (n Number) FloatVal() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Represents an entire document structure composed of keys and dynamo value
type Document map[string]interface{}

func (d Document) MarshalJSON() ([]byte, error) {
	output := make(map[string]interface{}, len(d))
	for key, val := range d {
		output[key] = wireEncode(val)
	}
	return json.Marshal(output)
}

func (d *Document) UnmarshalJSON(buf []byte) error {
	raw := make(map[string]interface{})
	err := json.Unmarshal(buf, &raw)
	if err != nil {
		return err
	}
	if *d == nil {
		*d = make(Document)
	}
	dd := *d

	for key, val := range raw {
		dd[key] = wireDecode(val)
	}
	return nil
}

// Helper to get a string from a document.
func (d Document) GetString(key string) string {
	return d[key].(string)
}

func (d Document) GetNumber(key string) Number {
	return d[key].(Number)
}

// Helper to build a hash key.
func HashKey(name string, value interface{}) Document {
	return Document{name: value}
}

// Helper to build a hash-range key.
func HashRangeKey(hashName string, hashVal interface{}, rangeName string, rangeVal interface{}) Document {
	return Document{
		hashName:  hashVal,
		rangeName: rangeVal,
	}
}
