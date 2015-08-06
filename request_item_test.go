package dynago

import (
	"testing"
)

func TestDeleteItem(t *testing.T) {
	assert, client, mock := setUp(t)
	di := client.DeleteItem("table", HashKey("Id", 50)).
		ConditionExpression("Foo = :bar", Param{":bar", "baz"}).
		ReturnValues(ReturnAllOld)
	assert.Equal(Document{"Id": 50}, di.req.Key)
	assert.Equal("Foo = :bar", di.req.ConditionExpression)
	assert.Equal(Document{":bar": "baz"}, di.req.ExpressionAttributeValues)

	attrib := Document{"Id": 50, "Foo": "Bar"}
	mock.DeleteItemResult = &DeleteItemResult{attrib}
	result, err := di.Execute()
	assert.NoError(err)
	assert.Equal(attrib, result.Attributes)
}

func TestGetItem(t *testing.T) {
	assert, client, mock := setUp(t)
	gi := client.GetItem("table", HashKey("Id", 10)).
		ProjectionExpression("#Foo, #Bar").
		Params(Param{"#Foo", "Foo"}).
		Param("#Bar", "BAR")
	assert.Equal(Document{"Id": 10}, gi.req.Key)
	assert.Equal("#Foo, #Bar", gi.req.ProjectionExpression)
	assert.Equal(map[string]string{"#Foo": "Foo", "#Bar": "BAR"}, gi.req.ExpressionAttributeNames)

	resultItem := Document{"ID": 10, "Name": "Foo"}
	mock.GetItemResult = &GetItemResult{Item: resultItem}
	result, err := gi.Execute()
	assert.NoError(err)
	assert.Equal(resultItem, result.Item)
}

func TestPutItem(t *testing.T) {
	assert, client, mock := setUp(t)
	doc := Document{"Id": 5}
	pi := client.PutItem("table", doc).
		ConditionExpression("#Foo = :Foo").
		Params(Document{"#Foo": "Foo", ":Foo": 50})
	assert.Equal(Document{":Foo": 50}, pi.req.ExpressionAttributeValues)
	response, err := pi.Execute()
	assert.NoError(err)
	assert.Nil(response)
	mock.PutItemResult = &PutItemResult{}
	pi = pi.ReturnValues(ReturnAllOld)
	response, err = pi.Execute()
	assert.NotNil(response)
}

// This test ensures that queries are properly copied (so they can be used across goroutines)
func TestQueryCopyProperty(t *testing.T) {
	assert, _, _ := setUp(t)
	q := &Query{}
	q2 := q.FilterExpression("hello")
	assert.Equal("hello", q2.req.FilterExpression)
	assert.Equal("", q.req.FilterExpression)
	q = q2

	assert.Nil(q.req.ConsistentRead)
	q.ConsistentRead(true)
	assert.Nil(q.req.ConsistentRead)

	assert.Nil(q.req.ScanIndexForward)
	q2 = q.Desc()
	assert.Nil(q.req.ScanIndexForward)
	assert.Equal(false, *q2.req.ScanIndexForward)

	q3 := q2.ScanIndexForward(true)
	assert.Nil(q.req.ScanIndexForward)
	assert.Equal(false, *q2.req.ScanIndexForward)
	assert.Equal(true, *q3.req.ScanIndexForward)
}

// This test checks that queries can be reused
func TestQueryReuse(t *testing.T) {
	assert, client, _ := setUp(t)
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

	q4 := q2.Params(Document{":val7": 95, ":offset": 8})
	check(q4, 3, 8)
	check(q2, 2, 45) // check clobbering again
}

func TestUpdateItem(t *testing.T) {
	assert, client, mock := setUp(t)
	ui := client.UpdateItem("table", HashKey("Id", 1)).
		ConditionExpression("Foo > :Foo").
		Params(Param{":Foo", "Bar"})
	assert.Equal(Document{":Foo": "Bar"}, ui.req.ExpressionAttributeValues)
	mock.UpdateItemResult = &UpdateItemResult{}
	result, err := ui.Execute()
	assert.NoError(err)
	assert.NotNil(result)
}
