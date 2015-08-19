package dynago

import (
	"testing"
)

func TestBatchGetItem(t *testing.T) {
	assert, client, mock := setUp(t)
	k1, k2, k3 := HashKey("Id", 5), HashKey("Id", 7), HashKey("Id", 10)
	bg := client.BatchGet().
		Get("table1", k1, k2).
		ProjectionExpression("table1", "Foo, #bar", Param{"#bar", "Bar"}).
		ConsistentRead("table1", true).
		Get("table2", k2, k3).
		ConsistentRead("table2", true).
		Get("table3", k1)

	tm := bg.buildTableMap()
	assert.Equal(3, len(tm))
	assert.Equal([]Document{k2, k1}, tm["table1"].Keys)
	assert.Equal("Foo, #bar", tm["table1"].ProjectionExpression)
	assert.Equal(map[string]string{"#bar": "Bar"}, tm["table1"].ExpressionAttributeNames)
	assert.Equal(true, tm["table1"].ConsistentRead)

	assert.Equal([]Document{k3, k2}, tm["table2"].Keys)
	assert.Equal("", tm["table2"].ProjectionExpression)
	assert.Equal(0, len(tm["table2"].ExpressionAttributeNames))
	assert.Equal(true, tm["table2"].ConsistentRead)

	assert.Equal(false, tm["table3"].ConsistentRead)
	assert.Equal([]Document{k1}, tm["table3"].Keys)

	mock.BatchGetItemResult = &BatchGetResult{}
	result, err := bg.Execute()
	assert.NoError(err)
	assert.NotNil(result)
}
