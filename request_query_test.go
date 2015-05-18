package dynago

import (
	"testing"
)

// This test ensures that queries are properly copied (so they can be used across goroutines)
func TestQueryCopyProperty(t *testing.T) {
	assert, _ := setUp(t)
	q := &Query{}
	q2 := q.FilterExpression("hello")
	assert.Equal("hello", q2.req.FilterExpression)
	assert.Equal("", q.req.FilterExpression)
	q = q2

	assert.Nil(q.req.ConsistentRead)
	q.ConsistentRead(true)
	assert.Nil(q.req.ConsistentRead)

	// TODO expand on this
}

// This test checks that queries can be reused
func TestQueryReuse(t *testing.T) {
	assert, client := setUp(t)
	q := client.Query("foo").
		ConsistentRead(true).
		KeyConditionExpression("Foo > :offset").
		Param(":unrelated", "hello")

	check := func(query *Query, length, val int) {
		assert.Equal(true, *query.req.ConsistentRead)
		assert.Equal("Foo > :offset", query.req.KeyConditionExpression)
		assert.Equal(length, len(query.req.ExpressionAttributeValues))
		assert.Equal(val, query.req.ExpressionAttributeValues[":offset"])
	}

	assert.Equal(1, len(q.req.ExpressionAttributeValues))
	q2 := q.Param(":offset", 45)
	check(q2, 2, 45)

	q3 := q.Param(":offset", 100)
	check(q3, 2, 100)
	check(q2, 2, 45) // Check we didn't clobber q2

	q4 := q2.Param(":val7", 95).Param(":offset", 8)
	check(q4, 3, 8)
	check(q2, 2, 45) // check clobbering again
}
