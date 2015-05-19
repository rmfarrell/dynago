package dynago

import (
	"strings"
)

type expressionAttributes struct {
	ExpressionAttributeValues Document          `json:",omitempty"`
	ExpressionAttributeNames  map[string]string `json:",omitempty"`
}

// Helper for a variety of endpoint types to build a params dictionary.
func (e *expressionAttributes) paramHelper(key string, value interface{}) {
	if strings.HasPrefix(key, "#") {
		if e.ExpressionAttributeNames == nil {
			e.ExpressionAttributeNames = map[string]string{key: value.(string)}
		} else {
			e.ExpressionAttributeNames[key] = value.(string)
		}
	} else {
		params := paramCopy(&e.ExpressionAttributeValues, 1)
		params[key] = value
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
