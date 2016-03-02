package dynago

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

/*
BinarySet stores a set of binary blobs in dynamo.

While implemented as a list in Go, DynamoDB does not preserve ordering on set
types and so may come back in a different order on retrieval. Use dynago.List
if ordering is important.
*/
type BinarySet [][]byte

/*
List represents DynamoDB lists, which are functionally very similar to JSON
lists.  Like JSON lists, these lists are heterogeneous, which means that the
elements of the list can be any valid value type, which includes other lists,
documents, numbers, strings, etc.
*/
type List []interface{}

/*
Return a copy of this list with all elements coerced as Documents.

It's very common to use lists in dynago where all elements in the list are
a Document. For that reason, this method is provided as a convenience to
get back your list as a list of documents.

If any element in the List is not a document, this will error.
As a convenience, even when it errors, a slice containing any elements
preceding the one which errored as documents will be given.
*/
func (l List) AsDocumentList() ([]Document, error) {
	docs := make([]Document, len(l))
	for i, listItem := range l {
		if doc, ok := listItem.(Document); !ok {
			return docs[:i], fmt.Errorf("item at index %d was not a Document", i)
		} else {
			docs[i] = doc
		}
	}
	return docs, nil
}

/*
A Number.

DynamoDB returns numbers as a string representation because they have a single
high-precision number type that can take the place of integers, floats, and
decimals for the majority of types.

This method has helpers to get the value of this number as one of various
Golang numeric types.
*/
type Number string

// IntVal interprets this number as an integer in base-10.
// error is returned if this is not a valid number or is too large.
func (n Number) IntVal() (int, error) {
	return strconv.Atoi(string(n))
}

// Int64Val interprets this number as an integer in base-10.
// error is returned if this string cannot be parsed as base 10 or is too large.
func (n Number) Int64Val() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// Uint64Val interprets this number as an unsigned integer.
// error is returned if this is not a valid positive integer or cannot fit.
func (n Number) Uint64Val() (uint64, error) {
	return strconv.ParseUint(string(n), 10, 64)
}

// FloatVal interprets this number as a floating point.
// error is returned if this number is not well-formed.
func (n Number) FloatVal() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

/*
NumberSet is an un-ordered set of numbers.

Sets in DynamoDB do not guarantee any ordering, so storing and retrieving a
NumberSet may not give you back the same order you put it in. The main
advantage of using sets in DynamoDB is using atomic updates with ADD and DELETE
in your UpdateExpression.
*/
type NumberSet []string

/*
Document is the core type for many dynamo operations on documents.
It is used to represent the root-level document, maps values, and
can also be used to supply expression parameters to queries.
*/
type Document map[string]interface{}

// MarshalJSON is used for encoding Document into wire representation.
func (d Document) MarshalJSON() ([]byte, error) {
	output := make(map[string]interface{}, len(d))
	for key, val := range d {
		if v := reflect.ValueOf(val); !isEmptyValue(v) {
			output[key] = wireEncode(val)
		}
	}
	return json.Marshal(output)
}

// UnmarshalJSON is used for unmarshaling Document from the wire representation.
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

/*
GetList gets the value at key as a List.

If value at key is nil, returns a nil list.
If value at key is not a List, will panic.
*/
func (d Document) GetList(key string) List {
	if d[key] != nil {
		return d[key].(List)
	} else {
		return nil
	}
}

/*
GetString gets the value at key as a String.

If the value at key is nil, returns an empty string.
If the value at key is not nil or a string, will panic.
*/
func (d Document) GetString(key string) string {
	if d[key] != nil {
		return d[key].(string)
	} else {
		return ""
	}
}

/*
GetNumber gets the value at key as a Number.

If the value at key is nil, returns a number containing the empty string.
If the value is not a Number or nil, will panic.
*/
func (d Document) GetNumber(key string) Number {
	if d[key] != nil {
		return d[key].(Number)
	} else {
		return Number("")
	}
}

/*
GetStringSet gets the value specified by key a StringSet.

If value at key is nil; returns an empty StringSet.
If it exists but is not a StringSet, panics.
*/
func (d Document) GetStringSet(key string) StringSet {
	if d[key] != nil {
		return d[key].(StringSet)
	} else {
		return StringSet{}
	}
}

/*
Helper to get a Time from a document.

If the value is omitted from the DB, or an empty string, then the return
is nil. If the value fails to parse as iso8601, then this method panics.
*/
func (d Document) GetTime(key string) (t *time.Time) {
	val := d[key]
	if val != nil {
		s := val.(string)
		parsed, err := time.ParseInLocation(iso8601compact, s, time.UTC)
		if err != nil {
			panic(err)
		}
		t = &parsed
	}
	return t
}

// AsParams makes Document satisfy the Params interface.
func (d Document) AsParams() (params []Param) {
	for key, val := range d {
		params = append(params, Param{key, val})
	}
	return
}

/*
GetBool gets the value specified by key as a boolean.

If the value does not exist in this Document, returns false.
If the value is the nil interface, also returns false.
If the value is a bool, returns the value of the bool.
If the value is a Number, returns true if value is non-zero.
For any other values, panics.
*/
func (d Document) GetBool(key string) bool {
	if v := d[key]; v != nil {
		switch val := v.(type) {
		case bool:
			return val
		case Number:
			if res, err := val.IntVal(); err != nil {
				panic(err)
			} else if res == 0 {
				return false
			} else {
				return true
			}
		default:
			panic(v)
		}
	} else {
		return false
	}
}

// HashKey is a shortcut to building keys used for various item operations.
func HashKey(name string, value interface{}) Document {
	return Document{name: value}
}

// HashRangeKey is a shortcut for building keys used for various item operations.
func HashRangeKey(hashName string, hashVal interface{}, rangeName string, rangeVal interface{}) Document {
	return Document{
		hashName:  hashVal,
		rangeName: rangeVal,
	}
}

// Param can be used as a single parameter.
type Param struct {
	Key   string
	Value interface{}
}

// AsParsms allows a solo Param to also satisfy the Params interface
func (p Param) AsParams() []Param {
	return []Param{p}
}

// P is a shortcut to create a single dynago Param.
// This is mainly for brevity especially with cross-package imports.
func P(key string, value interface{}) Params {
	return Param{key, value}
}

/*
Params encapsulates anything which can be used as expression parameters for
dynamodb expressions.

DynamoDB item queries using expressions can be provided parameters in a number
of handy ways:
	.Param(":k1", v1).Param(":k2", v2)
	-or-
	.Params(P(":k1", v1), P(":k2", v2))
	-or-
	.FilterExpression("...", P(":k1", v1), P(":k2", v2))
	-or-
	.FilterExpression("...", Document{":k1": v1, ":k2": v2})
Or any combination of Param, Document, or potentially other custom types which
provide the Params interface.
*/
type Params interface {
	AsParams() []Param
}

/*
StringSet is an un-ordered collection of distinct strings.

Sets in DynamoDB do not guarantee any ordering, so storing and retrieving a
StringSet may not give you back the same order you put it in. The main
advantage of using sets in DynamoDB is using atomic updates with ADD and DELETE
in your UpdateExpression.
*/
type StringSet []string

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}
