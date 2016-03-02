package dynago

import (
	"gopkg.in/underarmour/dynago.v1/schema"
)

/*
NewAwsClient is a shortcut to create a new dynamo client set up for AWS executor.

region is the AWS region, e.g. us-east-1.
accessKey is your amazon access key ID.
secretKey is your amazon secret key ID.
*/
func NewAwsClient(region string, accessKey string, secretKey string) *Client {
	endpoint := "https://dynamodb." + region + ".amazonaws.com/"
	return NewClient(NewAwsExecutor(endpoint, region, accessKey, secretKey))
}

/*
NewClient creates a new client from an executor.

For most use cases other than testing and mocking, you should be able to use
NewAwsClient which is a shortcut for this
*/
func NewClient(executor Executor) *Client {
	return &Client{executor, executor.SchemaExecutor()}
}

/*
Client is the primary start point of interaction with Dynago.

Client is concurrency safe, and completely okay to be used in multiple
threads/goroutines, as are the operations involving chaining on the client.
*/
type Client struct {
	executor       Executor
	schemaExecutor SchemaExecutor
}

/*
A BatchGet allows you to get up to 100 keys, in parallel, even across multiple
tables, in a single operation.
*/
func (c *Client) BatchGet() *BatchGet {
	return &BatchGet{client: c}
}

/*
A BatchWrite can compose a number of put or delete, even across multiple tables,
in a single operation.
*/
func (c *Client) BatchWrite() *BatchWrite {
	return newBatchWrite(c)
}

/*
DeleteItem lets you delete a single item in a table, optionally with conditions.
*/
func (c *Client) DeleteItem(table string, key Document) *DeleteItem {
	return newDeleteItem(c, table, key)
}

/*
GetItem gets a single document from a dynamo table.

key should be a Document containing enough attributes to describe the primary key.

You can use the HashKey or HashRangeKey helpers to help build a key:

	client.GetItem("foo", dynago.HashKey("Id", 45))

	client.GetItem("foo", dynago.HashRangeKey("UserId", 45, "Date", "20140512"))
*/
func (c *Client) GetItem(table string, key Document) *GetItem {
	return newGetItem(c, table, key)
}

/*
Query returns one or more items with range keys inside a single hash key.

This returns a new Query struct which you can compose via chaining to build the query you want.
Then finish the chain by calling Execute() to run the query.
*/
func (c *Client) Query(table string) *Query {
	return newQuery(c, table)
}

/*
PutItem creates or replaces a single document, optionally with conditions.

item should be a document representing the record and containing the attributes for the primary key.

Like all the other requests, you must call `Execute()` to run this.
*/
func (c *Client) PutItem(table string, item Document) *PutItem {
	return newPutItem(c, table, item)
}

/*
Scan can be used to enumerate an entire table or index.
*/
func (c *Client) Scan(table string) *Scan {
	return newScan(c, table)
}

/*
UpdateItem creates or modifies a single document.

UpdateItem will perform its updates even if the document at that key doesn't
exist; use a condition if you only want to update existing. documents.

The primary difference between PutItem and UpdateItem is the ability for
UpdateItem to use expressions for partial updates of a value.
*/
func (c *Client) UpdateItem(table string, key Document) *UpdateItem {
	return newUpdateItem(c, table, key)
}

/*
CreateTable makes a new table in your account.

This is not a synchronous operation; the table may take some time before
it is actually usable.
*/
func (c *Client) CreateTable(req *schema.CreateRequest) (*schema.CreateResult, error) {
	return c.schemaExecutor.CreateTable(req)
}

// DeleteTable deletes an existing table.
func (c *Client) DeleteTable(table string) (*schema.DeleteResult, error) {
	return c.schemaExecutor.DeleteTable(&schema.DeleteRequest{TableName: table})
}

// DescribeTable is used to get various attributes about the table.
func (c *Client) DescribeTable(table string) (*schema.DescribeResponse, error) {
	return c.schemaExecutor.DescribeTable(&schema.DescribeRequest{TableName: table})
}

// ListTables paginates through all the tables in an account.
func (c *Client) ListTables() *ListTables {
	return &ListTables{client: c}
}
