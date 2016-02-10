package dynago

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

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
	binaries := [][]byte{{'A', 'B'}, {'C', 'D'}}
	check(wireBinarySet{binaries}, BinarySet(binaries), `{"BS":["QUI=","Q0Q="]}`)

	// Nils
	check(wireNull{true}, nil, `{"NULL":true}`)

	// Numbers
	check(wireNumber{"7"}, int(7), `{"N":"7"}`)
	check(wireNumber{"-45"}, int64(-45), `{"N":"-45"}`)
	check(wireNumber{"4.55"}, float64(4.55), `{"N":"4.55"}`)
	check(wireNumberSet{[]string{"4", "5"}}, NumberSet{"4", "5"}, `{"NS":["4","5"]}`)

	check(wireNumber{"4500"}, uint(4500), `{"N":"4500"}`)
	check(wireNumber{"4500"}, uint64(4500), `{"N":"4500"}`)
	check(wireNumber{"4500"}, uint32(4500), `{"N":"4500"}`)
	check(wireNumber{"4500"}, uint16(4500), `{"N":"4500"}`)
	check(wireNumber{"251"}, uint8(251), `{"N":"251"}`)

	check(wireNumber{"123"}, int(123), `{"N":"123"}`)
	check(wireNumber{"123"}, int64(123), `{"N":"123"}`)
	check(wireNumber{"123"}, int32(123), `{"N":"123"}`)
	check(wireNumber{"123"}, int16(123), `{"N":"123"}`)
	check(wireNumber{"123"}, int8(123), `{"N":"123"}`)

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

	// Times
	time1 := time.Date(2014, 5, 5, 1, 2, 3, 0, time.UTC)
	Eastern, err := time.LoadLocation("US/Eastern")
	assert.NoError(err)
	check(wireString{"2014-05-05T01:02:03Z"}, time1, `{"S":"2014-05-05T01:02:03Z"}`)
	check(wireString{"2014-05-05T01:02:03Z"}, &time1, `{"S":"2014-05-05T01:02:03Z"}`)
	assert.Panics(func() { wireEncode(time1.In(Eastern)) })
}

func TestWireEncodeErrors(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		wireEncode([]int{1, 2})
	})
}

func TestWireDecode(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() { wireDecode(42) })
	decodeTest := func(k string, v interface{}) interface{} {
		return wireDecode(map[string]interface{}{k: v})
	}

	mapVal := map[string]interface{}{
		"Key1": map[string]interface{}{"S": "ABC"},
		"Key2": map[string]interface{}{"N": "123"},
	}
	listVal := []interface{}{mapVal["Key1"], mapVal["Key2"]}

	// Boolean
	assert.Equal(true, decodeTest("BOOL", true))
	assert.Equal(false, decodeTest("BOOL", false))
	// Binary
	assert.Equal([]byte("ABC"), decodeTest("B", "QUJD"))
	assert.Equal(BinarySet{[]byte("ABC"), []byte("AB")}, decodeTest("BS", []interface{}{"QUJD", "QUI="}))
	assert.Panics(func() { decodeTest("B", "QUJD=") })
	// Lists (heterogeneous)
	assert.Equal(List{"ABC", Number("123")}, decodeTest("L", listVal))
	// Maps (heterogeneous)
	assert.Equal(Document{"Key1": "ABC", "Key2": Number("123")}, decodeTest("M", mapVal))
	// Nil
	assert.Equal(nil, decodeTest("NULL", true))
	// Number
	assert.Equal(Number("45"), decodeTest("N", "45"))
	assert.Equal(NumberSet{"123", "456"}, decodeTest("NS", []interface{}{"123", "456"}))
	// Strings
	assert.Equal("FooBar", decodeTest("S", "FooBar"))
	assert.Equal(StringSet{"A", "B"}, decodeTest("SS", []interface{}{"A", "B"}))
}

func TestAnyInt(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() { anyInt("foo") })
	assert.Equal(int64(75), anyInt(int(75)))
	assert.Equal(int64(75), anyInt(int64(75)))
	assert.Equal(int64(75), anyInt(int32(75)))
	assert.Equal(int64(75), anyInt(int16(75)))
}
