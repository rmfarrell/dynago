package dynago

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const DynamoTargetPrefix = "DynamoDB_20120810." // This is the Dynamo API version we support

func NewClient(prefix string, region string, accessKey string, secretKey string) *Client {
	return &Client{
		endpoint: prefix,
		aws: AWSInfo{
			Region:    region,
			AccessKey: accessKey,
			SecretKey: secretKey,
			Service:   "dynamodb",
		},
	}
}

type Client struct {
	caller   http.Client
	endpoint string
	aws      AWSInfo
}

func (c *Client) makeRequest(target string, document interface{}) ([]byte, error) {
	buf, err := json.Marshal(document)
	log.Printf("Request Body: \n%s\n\n", buf)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", c.endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-amz-target", DynamoTargetPrefix+target)
	req.Header.Add("content-type", "application/x-amz-json-1.0")
	req.Header.Set("Host", req.URL.Host)
	c.aws.signRequest(req, buf)
	response, err := c.caller.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		resp := make([]byte, response.ContentLength)
		response.Body.Read(resp)
		log.Printf(" Another Status %d \n%#v\n\n%#v\n\n%s", response.StatusCode, req, response, resp)
	}
	return responseBytes(response)
}

func (c *Client) makeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error) {
	body, err := c.makeRequest(method, document)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, dest)
	return
}

/*
Return a Table object representing a table in dynamodb.
*/
func (c *Client) Table(tableName string) *Table {
	return &Table{c, tableName}
}

type Table struct {
	client    *Client
	tableName string
}

/*
Compose a GetItem on a dynamo table.

*/
func (c *Client) GetItem(table string, key Document) *GetItem {
	return newGetItem(c, table, key)
}

/*
Compose a Query on a dynamo table.

This returns a new Query struct which you can compose via chaining to build the query you want.

Then finish the chain by calling Execute() to run the query:

	result, err := table.Query().
		FilterExpression("Foo > :val").
		Param(":val", 45).
		Execute()

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

func responseBytes(response *http.Response) (buf []byte, err error) {
	if response.ContentLength > 0 {
		buf = make([]byte, response.ContentLength)
		_, err = response.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		err = response.Body.Close()
	}
	return
}
