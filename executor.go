package dynago

import (
	"encoding/json"

	"github.com/underarmour/dynago/internal/aws"
	"github.com/underarmour/dynago/schema"
)

/*
This interface defines how all the various queries manage their internal execution logic.

Executor is primarily provided so that testing and mocking can be done on
the API level, not just the transport level.

Executor can also optionally return a SchemaExecutor to execute schema actions.
*/
type Executor interface {
	BatchWriteItem(*BatchWrite) (*BatchWriteResult, error)
	GetItem(*GetItem) (*GetItemResult, error)
	PutItem(*PutItem) (*PutItemResult, error)
	Query(*Query) (*QueryResult, error)
	Scan(*Scan) (*ScanResult, error)
	UpdateItem(*UpdateItem) (*UpdateItemResult, error)
	SchemaExecutor() SchemaExecutor
}

type SchemaExecutor interface {
	CreateTable(*schema.CreateRequest) (*schema.CreateResult, error)
	DeleteTable(*schema.DeleteRequest) (*schema.DeleteResult, error)
	DescribeTable(*schema.DescribeRequest) (*schema.DescribeResponse, error)
	ListTables(*ListTables) (*schema.ListResponse, error)
}

type AwsRequester interface {
	MakeRequest(target string, body []byte) ([]byte, error)
}

type awsExecutor struct {
	Requester AwsRequester
}

// Create an AWS executor with a specified endpoint and AWS parameters.
func NewAwsExecutor(endpoint, region, accessKey, secretKey string) Executor {
	signer := aws.AwsInfo{
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Service:   "dynamodb",
	}
	requester := &aws.RequestMaker{
		Endpoint:       endpoint,
		Signer:         &signer,
		BuildError:     buildError,
		DebugRequests:  Debug.HasFlag(DebugRequests),
		DebugResponses: Debug.HasFlag(DebugResponses),
		DebugFunc:      DebugFunc,
	}
	return &awsExecutor{requester}
}

func (e *awsExecutor) makeRequest(target string, document interface{}) ([]byte, error) {
	buf, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}
	return e.requester.MakeRequest(target, buf)
}

func (e *awsExecutor) makeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error) {
	body, err := e.makeRequest(method, document)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, dest)
	return
}

func (e *awsExecutor) SchemaExecutor() SchemaExecutor {
	return e
}
