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
		Get("table2", k2, k3)

	tm := bg.buildTableMap()
	assert.Equal(2, len(tm))
	assert.Equal([]Document{k2, k1}, tm["table1"].Keys)
	assert.Equal("Foo, #bar", tm["table1"].ProjectionExpression)
	assert.Equal(map[string]string{"#bar": "Bar"}, tm["table1"].ExpressionAttributeNames)

	assert.Equal([]Document{k3, k2}, tm["table2"].Keys)
	assert.Equal("", tm["table2"].ProjectionExpression)
	assert.Equal(0, len(tm["table2"].ExpressionAttributeNames))

	mock.BatchGetItemResult = &BatchGetResult{}
	result, err := bg.Execute()
	assert.NoError(err)
	assert.NotNil(result)
}
