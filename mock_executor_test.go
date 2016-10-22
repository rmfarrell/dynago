package dynago_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rmfarrell/dynago"
)

func mockSetup(t *testing.T) (*assert.Assertions, *dynago.Client, *dynago.MockExecutor) {
	t.Parallel()
	executor := &dynago.MockExecutor{}
	return assert.New(t), dynago.NewClient(executor), executor
}

func TestMockExecutorBatchGetItem(t *testing.T) {
	assert, client, executor := mockSetup(t)

	key1, key2 := dynago.HashKey("Id", 3), dynago.HashKey("Id", 4)
	client.BatchGet().
		Get("table1", key1, key2).
		ProjectionExpression("table1", "#n,Foo,Bar", dynago.P("#n", "Name")).
		Get("table2", key2).
		Execute()
	assert.Equal(true, executor.BatchGetItemCalled)
	assert.NotNil(executor.BatchGetItemCall)
	call := executor.BatchGetItemCall
	assert.Equal("BatchGetItem", call.Method)
	assert.Equal(2, len(call.BatchGets))
	assert.Equal([]dynago.Document{key2, key1}, call.BatchGets["table1"].Keys)
	assert.Equal("#n,Foo,Bar", call.BatchGets["table1"].ProjectionExpression)
	assert.Equal(map[string]string{"#n": "Name"}, call.BatchGets["table1"].ExpressionAttributeNames)
	assert.Equal([]dynago.Document{key2}, call.BatchGets["table2"].Keys)
	assert.Equal("", call.BatchGets["table2"].ProjectionExpression)
}

func TestMockExecutorBatchWriteItem(t *testing.T) {
	assert, client, executor := mockSetup(t)
	doc1 := dynago.Document{"Id": 1, "Name": "1"}
	doc2 := dynago.Document{"Id": 2, "Name": "2"}
	key3 := dynago.HashKey("Id", 3)
	key4 := dynago.HashKey("Id", 4)
	client.BatchWrite().
		Put("table1", doc1, doc2).
		Delete("table1", key3).
		Delete("table2", key4).
		Execute()
	assert.Equal(true, executor.BatchWriteItemCalled)
	assert.Equal([]dynago.Document{doc2, doc1}, executor.BatchWriteItemCall.BatchWrites.GetPuts("table1"))
	assert.Equal(0, len(executor.BatchWriteItemCall.BatchWrites.GetPuts("table2")))
	assert.Equal([]dynago.Document{key3}, executor.BatchWriteItemCall.BatchWrites.GetDeleteKeys("table1"))
	assert.Equal([]dynago.Document{key4}, executor.BatchWriteItemCall.BatchWrites.GetDeleteKeys("table2"))
	assert.Equal(executor.Calls[0], *executor.BatchWriteItemCall)
}

func TestMockExecutorDeleteItem(t *testing.T) {
	assert, client, executor := mockSetup(t)
	client.DeleteItem("table1", dynago.HashKey("Id", 51)).
		ConditionExpression("expr1", dynago.P(":foo", 4), dynago.P("#f", "f")).
		ReturnValues(dynago.ReturnAllOld).
		Execute()
	assert.Equal(true, executor.DeleteItemCalled)
	assert.Equal(executor.Calls[0], *executor.DeleteItemCall)
	call := executor.DeleteItemCall
	assert.Equal("DeleteItem", call.Method)
	assert.Equal(dynago.HashKey("Id", 51), call.Key)
	assert.Equal("table1", call.Table)
	assert.Equal("expr1", call.ConditionExpression)
	assert.Equal(dynago.ReturnAllOld, call.ReturnValues)
	assert.Equal(dynago.Document{":foo": 4}, call.ExpressionAttributeValues)
	assert.Equal(map[string]string{"#f": "f"}, call.ExpressionAttributeNames)
}

func TestMockExecutorGetItem(t *testing.T) {
	assert, client, executor := mockSetup(t)
	client.GetItem("table1", dynago.HashKey("Id", 5)).
		ConsistentRead(true).
		ProjectionExpression("foo").
		Param(":foo", "bar").
		Execute()

	assert.Equal(true, executor.GetItemCalled)
	assert.Equal("GetItem", executor.GetItemCall.Method)
	assert.Equal("table1", executor.GetItemCall.Table)
	assert.Equal(true, executor.GetItemCall.ConsistentRead)
	assert.Equal(dynago.Document{"Id": 5}, executor.GetItemCall.Key)
	assert.Equal("foo", executor.GetItemCall.ProjectionExpression)
	assert.Equal(dynago.Document{":foo": "bar"}, executor.GetItemCall.ExpressionAttributeValues)
	assert.Equal(1, len(executor.Calls))
	assert.Equal(executor.Calls[0], *executor.GetItemCall)
}

func TestMockExecutorPutItem(t *testing.T) {
	assert, client, executor := mockSetup(t)
	client.PutItem("table2", dynago.HashKey("Id", 5)).
		ConditionExpression("Foo = :bar").Param(":bar", "45").
		ReturnValues(dynago.ReturnUpdatedNew).
		Execute()
	assert.Equal(true, executor.PutItemCalled)
	assert.Equal("PutItem", executor.PutItemCall.Method)
	assert.Equal("table2", executor.PutItemCall.Table)
	assert.Equal(dynago.HashKey("Id", 5), executor.PutItemCall.Item)
	assert.Equal(dynago.ReturnUpdatedNew, executor.PutItemCall.ReturnValues)
	assert.Equal(executor.Calls[0], *executor.PutItemCall)
}

func TestMockExecutorQuery(t *testing.T) {
	assert, client, executor := mockSetup(t)
	client.Query("table3").IndexName("Index1").
		ConsistentRead(true).
		KeyConditionExpression("ABC = :def").
		FilterExpression("Foo > :param", dynago.P(":param", 95)).
		Limit(50).Select(dynago.SelectSpecificAttributes).
		Execute()
	assert.Equal(true, executor.QueryCalled)
	assert.Equal("Query", executor.QueryCall.Method)
	assert.Equal("table3", executor.QueryCall.Table)
	assert.Equal("Index1", executor.QueryCall.IndexName)
	assert.Equal(true, executor.QueryCall.ConsistentRead)
	assert.Equal(true, executor.QueryCall.Ascending)
	assert.Equal(uint(50), executor.QueryCall.Limit)
	assert.Equal(dynago.Document{":param": 95}, executor.QueryCall.ExpressionAttributeValues)
	assert.Equal("ABC = :def", executor.QueryCall.KeyConditionExpression)
	assert.Equal("Foo > :param", executor.QueryCall.FilterExpression)
	assert.Equal(dynago.SelectSpecificAttributes, executor.QueryCall.Select)
	assert.Equal(executor.Calls[0], *executor.QueryCall)

	doc1 := dynago.Document{"Id": 1, "Name": "1"}
	executor.QueryResult = &dynago.QueryResult{Items: []dynago.Document{doc1}}
	executor.QueryError = &dynago.Error{}
	result, err := client.Query("table3").Desc().Execute()
	assert.Equal(2, len(executor.Calls))
	assert.Equal(false, executor.QueryCall.Ascending)
	assert.Empty(false, executor.QueryCall.ConsistentRead)
	assert.Equal(executor.Calls[1], *executor.QueryCall)
	assert.Error(err)
	assert.NotNil(result)
	assert.Equal(1, len(result.Items))
	assert.Equal(1, result.Count)
	assert.Equal(1, result.ScannedCount)
}

func TestMockExecutorScan(t *testing.T) {
	assert, client, executor := mockSetup(t)
	scan := client.Scan("table5").
		ExclusiveStartKey(dynago.HashKey("Id", 2)).
		FilterExpression("Foo = :bar", dynago.P(":bar", 10)).
		ProjectionExpression("Foo, Bar, #baz", dynago.P("#baz", "Baz")).
		IndexName("index5")
	scan.Execute()
	assert.Equal(true, executor.ScanCalled)
	assert.NotNil(executor.ScanCall)
	assert.Equal("Scan", executor.ScanCall.Method)
	assert.Equal("table5", executor.ScanCall.Table)
	assert.Equal("Foo = :bar", executor.ScanCall.FilterExpression)
	assert.Equal(dynago.Document{":bar": 10}, executor.ScanCall.ExpressionAttributeValues)
	assert.Equal(map[string]string{"#baz": "Baz"}, executor.ScanCall.ExpressionAttributeNames)
	assert.Equal("Foo, Bar, #baz", executor.ScanCall.ProjectionExpression)
	assert.Equal("index5", executor.ScanCall.IndexName)
	assert.Nil(executor.ScanCall.Segment)
	scan.Segment(5, 10).Select(dynago.SelectCount).Execute()
	assert.Equal(2, len(executor.Calls))
	assert.NotNil(executor.ScanCall.Segment)
	assert.Equal(5, *executor.ScanCall.Segment)
	assert.Equal(10, *executor.ScanCall.TotalSegments)
	assert.Equal(dynago.SelectCount, executor.ScanCall.Select)
}

func TestMockExecutorUpdateItem(t *testing.T) {
	assert, client, executor := mockSetup(t)

	client.UpdateItem("table4", dynago.HashKey("Id", 50)).
		UpdateExpression("Foo = :param1").Param(":param1", 90).
		ConditionExpression("#foo > :param2").Param("#foo", "Foo").
		Execute()
	assert.Equal(true, executor.UpdateItemCalled)
	assert.Equal("UpdateItem", executor.UpdateItemCall.Method)
	assert.Equal("table4", executor.UpdateItemCall.Table)
	assert.Equal("Foo = :param1", executor.UpdateItemCall.UpdateExpression)
	assert.Equal("#foo > :param2", executor.UpdateItemCall.ConditionExpression)
	assert.Equal(map[string]string{"#foo": "Foo"}, executor.UpdateItemCall.ExpressionAttributeNames)
	assert.Equal(dynago.Document{":param1": 90}, executor.UpdateItemCall.ExpressionAttributeValues)
}
