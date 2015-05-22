package dynago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setUp(t *testing.T) (*assert.Assertions, *Client, *MockExecutor) {
	executor := &MockExecutor{}
	return assert.New(t), NewClient(executor), executor
}

func TestQueryParams(t *testing.T) {
	assert, client, _ := setUp(t)
	q := client.Query("Foo").Param(":start", 9)
	assert.Equal(1, len(q.req.ExpressionAttributeValues))
	q2 := q.Params(Param{":end", 4}, Param{":other", "hello"}, Param{"#name", "Name"})
	assert.Equal(3, len(q2.req.ExpressionAttributeValues))
	assert.Equal(1, len(q2.req.ExpressionAttributeNames))
	assert.Equal("Name", q2.req.ExpressionAttributeNames["#name"])
}
