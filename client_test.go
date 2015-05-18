package dynago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setUp(t *testing.T) (*assert.Assertions, *Client) {
	// TODO add the mock executor
	return assert.New(t), NewClientExecutor(nil)
}
