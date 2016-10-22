package dynago

import (
	"encoding/json"

	"github.com/rmfarrell/dynago/internal/aws"
	"github.com/rmfarrell/dynago/schema"
)

/*
Executor defines how all the various queries manage their internal execution logic.

Executor is primarily provided so that testing and mocking can be done on
the API level, not just the transport level.

Executor can also optionally return a SchemaExecutor to execute schema actions.
*/
type Executor interface {
	BatchGetItem(*BatchGet) (*BatchGetResult, error)
	BatchWriteItem(*BatchWrite) (*BatchWriteResult, error)
	DeleteItem(*DeleteItem) (*DeleteItemResult, error)
	GetItem(*GetItem) (*GetItemResult, error)
	PutItem(*PutItem) (*PutItemResult, error)
	Query(*Query) (*QueryResult, error)
	Scan(*Scan) (*ScanResult, error)
	UpdateItem(*UpdateItem) (*UpdateItemResult, error)
	SchemaExecutor() SchemaExecutor
}

// SchemaExecutor implements schema management commands.
type SchemaExecutor interface {
	CreateTable(*schema.CreateRequest) (*schema.CreateResult, error)
	DeleteTable(*schema.DeleteRequest) (*schema.DeleteResult, error)
	DescribeTable(*schema.DescribeRequest) (*schema.DescribeResponse, error)
	ListTables(*ListTables) (*schema.ListResponse, error)
}

// AwsRequester makes requests to dynamodb
type AwsRequester interface {
	MakeRequest(target string, body []byte) ([]byte, error)
}

// Create an AWS executor with a specified endpoint and AWS parameters.
func NewAwsExecutor(endpoint, region, accessKey, secretKey string) *AwsExecutor {
	signer := aws.AwsSigner{
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Service:   "dynamodb",
	}
	requester := &aws.RequestMaker{
		Endpoint:       aws.FixEndpointUrl(endpoint),
		Signer:         &signer,
		BuildError:     buildError,
		DebugRequests:  Debug.HasFlag(DebugRequests),
		DebugResponses: Debug.HasFlag(DebugResponses),
		DebugFunc:      DebugFunc,
	}
	return &AwsExecutor{requester}
}

/*
AwsExecutor is the underlying implementation of making requests to DynamoDB.
*/
type AwsExecutor struct {
	// Underlying implementation that makes requests for this executor. It
	// is called to make every request that the executor makes. Swapping the
	// underlying implementation is not thread-safe and therefore not
	// recommended in production code.
	Requester AwsRequester
}

func (e *AwsExecutor) makeRequest(target string, document interface{}) ([]byte, error) {
	buf, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}
	return e.Requester.MakeRequest(target, buf)
}

/*
Make a request to the underlying requester, marshaling document as JSON,
and if the requester doesn't error, unmarshaling the response back into dest.

This method is mostly exposed for those implementing custom executors or
prototyping new functionality.
*/
func (e *AwsExecutor) MakeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error) {
	body, err := e.makeRequest(method, document)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, dest)
	return
}

// Return a SchemaExecutor making requests on this Executor.
func (e *AwsExecutor) SchemaExecutor() SchemaExecutor {
	return awsSchemaExecutor{e}
}
