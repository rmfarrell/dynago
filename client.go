package dynago

import (
	"github.com/underarmour/dynago/schema"
)

/*
Create a new dynamo client set up for AWS executor.

region is the AWS region, e.g. us-east-1.
accessKey is your amazon access key ID.
secretKey is your amazon secret key ID.
*/
func NewAwsClient(region string, accessKey string, secretKey string) *Client {
	endpoint := "https://dynamodb." + region + ".amazonaws.com/"
	return NewClient(NewAwsExecutor(endpoint, region, accessKey, secretKey))
}

/*
Create a new client.

For most use cases other than testing and mocking, you should be able to use
NewAwsClient which is a shortcut for this
*/
func NewClient(executor Executor) *Client {
	return &Client{executor}
}

type Client struct {
	executor Executor
}

/*
Compose a batch write.

Batch writes can compose a number of put or delete, even across multiple tables, in a single operation.
*/
func (c *Client) BatchWrite() *BatchWrite {
	return newBatchWrite(c)
}

/*
Compose a GetItem on a dynamo table.

key should be a Document containing enough attributes to describe the primary key.

You can use the HashKey or HashRangeKey helpers to help build a key:

	client.GetItem("foo", dynago.HashKey("Id", 45))

	client.GetItem("foo", dynago.HashRangeKey("UserId", 45, "Date", "20140512"))
*/
func (c *Client) GetItem(table string, key Document) *GetItem {
	return newGetItem(c, table, key)
}

/*
Compose a Query on a dynamo table.

This returns a new Query struct which you can compose via chaining to build the query you want.
Then finish the chain by calling Execute() to run the query.
*/
func (c *Client) Query(table string) *Query {
	return newQuery(c, table)
}

/*
Compose a PutItem on a dynamo table.

item should be a document representing the record and containing the attributes for the primary key.

Like all the other requests, you must call `Execute()` to run this.
*/
func (c *Client) PutItem(table string, item Document) *PutItem {
	return newPutItem(c, table, item)
}

/*
Compose an UpdateItem on a dynamo table.
*/
func (c *Client) UpdateItem(table string, key Document) *UpdateItem {
	return newUpdateItem(c, table, key)
}

/*
Create a table.
*/
func (c *Client) CreateTable(req *schema.CreateRequest) (*schema.CreateResponse, error) {
	return c.executor.CreateTable(req)
}
