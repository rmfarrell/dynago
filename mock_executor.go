package dynago

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
		assert.Equal("mytable", executor.PutItemCall.Table)
		... and so on
	}

*/
type MockExecutor struct {
	Calls []MockExecutorCall // All calls made through this executor

	PutItemCalled bool
	PutItemCall   *MockExecutorCall
	PutItemResult *PutItemResult
	PutItemError  error

	GetItemCalled     bool
	GetItemCall       *MockExecutorCall
	GetItemResultItem Document
	GetItemError      error

	BatchWriteItemCalled bool
	BatchWriteItemCall   *MockExecutorCall
	BatchWriteItemError  error

	QueryCalled bool              // True if query was called at least once
	QueryCall   *MockExecutorCall // Info for the last call to Query
	QueryError  error             // Specify the error from Query
	QueryResult *QueryResult      // Specify the result from Query

	UpdateItemCalled bool
	UpdateItemCall   *MockExecutorCall
	UpdateItemResult *UpdateItemResult
	UpdateItemError  error
}

// Mock executor calls
type MockExecutorCall struct {
	// used for all calls
	Method string
	Table  string

	// used for calls with expressions (most of them)
	ExpressionAttributeNames  map[string]string
	ExpressionAttributeValues Document

	Key                 Document
	Item                Document
	UpdateExpression    string
	ConditionExpression string
	ReturnValues        ReturnValues
	ConsistentRead      bool

	// Query only
	IndexName              string
	KeyConditionExpression string
	FilterExpression       string
	Ascending              bool
	Limit                  uint
	ExclusiveStartKey      Document

	BatchWrites BatchWriteTableMap
}

func (e *MockExecutor) BatchWriteItem(batchWrite *BatchWrite) (*BatchWriteResult, error) {
	e.BatchWriteItemCalled = true
	e.addCall(&e.BatchWriteItemCall, MockExecutorCall{
		Method:      "BatchWriteItem",
		BatchWrites: batchWrite.buildTableMap(),
	})

	return &BatchWriteResult{}, e.BatchWriteItemError
}

func (e *MockExecutor) GetItem(getItem *GetItem) (*GetItemResult, error) {
	e.GetItemCalled = true
	call := MockExecutorCall{
		Method:         "GetItem",
		Table:          getItem.req.TableName,
		Key:            getItem.req.Key,
		ConsistentRead: getItem.req.ConsistentRead,
	}
	e.GetItemCall = &call
	e.Calls = append(e.Calls, call)
	return &GetItemResult{Item: e.GetItemResultItem}, e.GetItemError
}

func (e *MockExecutor) PutItem(putItem *PutItem) (*PutItemResult, error) {
	e.PutItemCalled = true
	call := MockExecutorCall{
		Method:                    "PutItem",
		Table:                     putItem.req.TableName,
		Item:                      putItem.req.Item,
		ReturnValues:              putItem.req.ReturnValues,
		ConditionExpression:       putItem.req.ConditionExpression,
		ExpressionAttributeNames:  putItem.req.ExpressionAttributeNames,
		ExpressionAttributeValues: putItem.req.ExpressionAttributeValues,
	}
	e.PutItemCall = &call
	e.Calls = append(e.Calls, call)
	return e.PutItemResult, e.PutItemError
}

func (e *MockExecutor) Query(query *Query) (*QueryResult, error) {
	e.QueryCalled = true
	ascending, consistent := true, false
	if query.req.ScanIndexForward != nil {
		ascending = *query.req.ScanIndexForward
	}
	if query.req.ConsistentRead != nil {
		consistent = *query.req.ConsistentRead
	}

	e.addCall(&e.QueryCall, MockExecutorCall{
		Method:                    "Query",
		Table:                     query.req.TableName,
		IndexName:                 query.req.IndexName,
		KeyConditionExpression:    query.req.KeyConditionExpression,
		FilterExpression:          query.req.FilterExpression,
		ExpressionAttributeNames:  query.req.ExpressionAttributeNames,
		ExpressionAttributeValues: query.req.ExpressionAttributeValues,
		Ascending:                 ascending,
		ConsistentRead:            consistent,
		Limit:                     query.req.Limit,
		ExclusiveStartKey:         query.req.ExclusiveStartKey,
	})

	result := e.QueryResult
	if result != nil {
		result.query = query
		if result.Count == 0 {
			result.Count = len(result.Items)
		}
		if result.ScannedCount == 0 {
			result.ScannedCount = result.Count
		}
	}

	return result, e.QueryError
}

func (e *MockExecutor) UpdateItem(update *UpdateItem) (*UpdateItemResult, error) {
	e.UpdateItemCalled = true
	e.addCall(&e.UpdateItemCall, MockExecutorCall{
		Method:                    "UpdateItem",
		Table:                     update.req.TableName,
		UpdateExpression:          update.req.UpdateExpression,
		ConditionExpression:       update.req.ConditionExpression,
		ExpressionAttributeNames:  update.req.ExpressionAttributeNames,
		ExpressionAttributeValues: update.req.ExpressionAttributeValues,
	})
	return e.UpdateItemResult, e.UpdateItemError
}

// Currently we don't implement mocking for SchemaExecutor. Returns nil.
func (e *MockExecutor) SchemaExecutor() SchemaExecutor {
	return nil
}

// Reduce boilerplate on adding a call
func (e *MockExecutor) addCall(target **MockExecutorCall, call MockExecutorCall) {
	e.Calls = append(e.Calls, call)
	if target != nil {
		*target = &call
	}
}

// Convenience method to get delete keys in this batch write map for a specific table
func (m BatchWriteTableMap) GetDeleteKeys(table string) (deletes []Document) {
	for _, entry := range m[table] {
		if entry.DeleteRequest != nil {
			deletes = append(deletes, entry.DeleteRequest.Key)
		}
	}
	return
}

// Convenience method to get all put documents in this batch write map for a specific table
func (m BatchWriteTableMap) GetPuts(table string) (puts []Document) {
	for _, entry := range m[table] {
		if entry.PutRequest != nil {
			puts = append(puts, entry.PutRequest.Item)
		}
	}
	return
}
