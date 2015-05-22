package dynago

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWireEncodeBasic(t *testing.T) {
	assert := assert.New(t)
	check := func(expected interface{}, value interface{}, expectedJson string) interface{} {
		encoded := wireEncode(value)
		r := reflect.ValueOf(encoded)
		assert.Equal(reflect.Ptr, r.Kind())
		deref := reflect.Indirect(r).Interface()
		assert.IsType(expected, deref)
		assert.Equal(expected, deref)
		if expectedJson != "" {
			b, err := json.Marshal(encoded)
			assert.NoError(err)
			assert.Equal(expectedJson, string(b))
		}
		return encoded
	}
	// Booleans
	check(wireBool{true}, true, `{"BOOL":true}`)
	check(wireBool{false}, false, `{"BOOL":false}`)

	// Binary
	check(wireBinary{[]byte{'A', 'B'}}, []byte{'A', 'B'}, `{"B":"QUI="}`)
	binaries := [][]byte{[]byte{'A', 'B'}, []byte{'C', 'D'}}
	check(wireBinarySet{binaries}, BinarySet(binaries), `{"BS":["QUI=","Q0Q="]}`)

	// Numbers
	check(wireNumber{"7"}, int(7), `{"N":"7"}`)
	check(wireNumber{"-45"}, int64(-45), `{"N":"-45"}`)
	check(wireNumber{"4.55"}, float64(4.55), `{"N":"4.55"}`)
	check(wireNumberSet{[]string{"4", "5"}}, NumberSet{"4", "5"}, `{"NS":["4","5"]}`)

	// Lists (heterogeneous)
	check(
		wireList{[]interface{}{&wireNumber{"45"}, &wireString{"Hello"}, &wireNumber{"4.5"}, &wireBinary{[]byte("AB")}}},
		List{45, "Hello", float64(4.5), []byte("AB")},
		`{"L":[{"N":"45"},{"S":"Hello"},{"N":"4.5"},{"B":"QUI="}]}`,
	)

	// Maps / Documents
	eMap := map[string]interface{}{"Foo": &wireNumber{"42"}}
	jsonMap := `{"M":{"Foo":{"N":"42"}}}`
	check(wireMap{eMap}, Document{"Foo": 42}, jsonMap)
	check(wireMap{eMap}, map[string]interface{}{"Foo": 42}, jsonMap)

	// Strings
	check(wireString{"Foo"}, "Foo", `{"S":"Foo"}`)
	check(wireStringSet{[]string{"A", "B"}}, StringSet{"A", "B"}, `{"SS":["A","B"]}`)

}

func TestWireEncodeErrors(t *testing.T) {
	// TODO
}
