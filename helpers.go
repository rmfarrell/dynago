package dynago

import (
	"fmt"
	"strings"
)

type expressionAttributes struct {
	ExpressionAttributeValues Document          `json:",omitempty"`
	ExpressionAttributeNames  map[string]string `json:",omitempty"`
}

// Helper for a variety of endpoint types to build a params dictionary.
func (e *expressionAttributes) paramHelper(key string, value interface{}) {
	e.assignParams([]Param{{key, value}})
}

// Helper to build multi-params dictionary
func (e *expressionAttributes) paramsHelper(params []interface{}) {
	if len(params) == 0 {
		return
	}
	output := make([]Param, 0, len(params))
	for _, p := range params {
		switch value := p.(type) {
		case Param:
			output = append(output, value)
		case *Param:
			output = append(output, *value)
		case Document:
			for k, v := range value {
				output = append(output, Param{k, v})
			}
		default:
			panic(fmt.Errorf("Don't know what to do with value %#v", p))
		}
	}
	e.assignParams(output)
}

func (e *expressionAttributes) assignParams(params []Param) {
	var copyValues, copyNames bool
	for _, p := range params {
		if strings.HasPrefix(p.Key, "#") {
			if !copyNames {
				copyNames = true
				eaNameCopy(&e.ExpressionAttributeNames, 1)
			}
			e.ExpressionAttributeNames[p.Key] = p.Value.(string)
		} else {
			if !copyValues {
				copyValues = true
				paramCopy(&e.ExpressionAttributeValues, len(params))
			}
			e.ExpressionAttributeValues[p.Key] = p.Value
		}
	}
}

func paramCopy(doc *Document, extendBy int) Document {
	params := make(Document, len(*doc)+extendBy)
	if *doc != nil {
		for k, v := range *doc {
			params[k] = v
		}
	}
	*doc = params
	return params
}

func eaNameCopy(doc *map[string]string, extendBy int) {
	names := make(map[string]string, len(*doc)+extendBy)
	for k, v := range *doc {
		names[k] = v
	}
	*doc = names
}

/*
Set the debug mode.

This is a set of bit-flags you can use to set up how much debugging dynago uses:
	dynago.Debug = dynago.DebugRequests | dynago.DebugResponses
*/
var Debug DebugFlags

// Set the target of debug. Must be set for debug to be used.
var DebugFunc func(format string, v ...interface{})
