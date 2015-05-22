package dynago

import "github.com/underarmour/dynago/schema"

/*
A Mock executor for purpose of testing.

This Executor doesn't actually run any network requests, so it can be used in
unit testing for your own application which uses Dynago.  It can be asserted
on whether a specific underlying method got called, what it was called with,
and the nature of how the query was built. You can also for most methods
control which result is returned in order to control how results make it back
to the application.

	// example, normally you'd call into a real application not inside your
	// test module and use dependency injection to specify the client.
	func application(client *dynago.Client) int {
		result, err := client.PutItem("mytable", doc).Execute()
		// do something with result maybe.
	}

	func TestApplication() {
		executor := &dynago.MockExecutor{}
		client := dynago.NewClient(executor)
		executor.PutItemResult = &dynago.PutItemResult{}

		// call into application
		application(client)

		// assert things on the executor.
		assert.Equal(true, executor.PutItemCalled)
		assert.Equal("mytable", executor.PutItemCalledWithTable)
		... and so on
	}

*/
type MockExecutor struct {
	PutItemCalled          bool
	PutItemCalledWithTable string
	PutItemCalledWithItem  Document
	PutItemResult          *PutItemResult
	PutItemError           error

	GetItemCalled          bool
	GetItemCalledWithTable string
	GetItemCalledWithKey   Document
	GetItemResultItem      Document
	GetItemError           error

	BatchWriteItemCalled                  bool
	BatchWriteItemCalledWithDeleteTable   string
	BatchWriteItemCalledWithDeleteActions []Document
	BatchWriteItemError                   error

	QueryCalled                              bool
	QueryCalledWithTable                     string
	QueryCalledWithIndexName                 string
	QueryCalledWithKeyConditionExpression    string
	QueryCalledWithExpressionAttributeValues Document
	QueryError                               error
	QueryResultItems                         []Document
}

func (e *MockExecutor) GetItem(getItem *GetItem) (*GetItemResult, error) {
	e.GetItemCalled = true
	e.GetItemCalledWithTable = getItem.req.TableName
	e.GetItemCalledWithKey = getItem.req.Key
	return &GetItemResult{Item: e.GetItemResultItem}, e.GetItemError
}

func (e *MockExecutor) PutItem(putItem *PutItem) (*PutItemResult, error) {
	e.PutItemCalled = true
	e.PutItemCalledWithTable = putItem.req.TableName
	e.PutItemCalledWithItem = putItem.req.Item
	return e.PutItemResult, e.PutItemError
}

func (e *MockExecutor) Query(query *Query) (*QueryResult, error) {
	e.QueryCalled = true
	e.QueryCalledWithTable = query.req.TableName
	e.QueryCalledWithIndexName = query.req.IndexName
	e.QueryCalledWithKeyConditionExpression = query.req.KeyConditionExpression
	e.QueryCalledWithExpressionAttributeValues = query.req.ExpressionAttributeValues
	return &QueryResult{Items: e.QueryResultItems}, e.QueryError
}

func (e *MockExecutor) UpdateItem(*UpdateItem) (*UpdateItemResult, error) {
	return nil, nil
}

func (e *MockExecutor) CreateTable(*schema.CreateRequest) (*schema.CreateResponse, error) {
	return nil, nil
}

func (e *MockExecutor) BatchWriteItem(batchWrite *BatchWrite) (*BatchWriteResult, error) {
	e.BatchWriteItemCalled = true
	e.BatchWriteItemCalledWithDeleteTable = batchWrite.deletes.table
	for deleteRequest := batchWrite.deletes; deleteRequest != nil; deleteRequest = batchWrite.deletes.next {
		e.BatchWriteItemCalledWithDeleteActions = append(e.BatchWriteItemCalledWithDeleteActions, deleteRequest.item)
	}
	return nil, e.BatchWriteItemError
}
