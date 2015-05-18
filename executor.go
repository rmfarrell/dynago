package dynago

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/crast/dynago/schema"
)

/*
This interface defines how all the various queries manage their internal execution logic.

Executor is primarily provided so that testing and mocking can be done on
the API level, not just the transport level.
*/
type Executor interface {
	GetItem(*GetItem) (*GetItemResult, error)
	PutItem(*PutItem) (*PutItemResult, error)
	Query(*Query) (*QueryResult, error)
	UpdateItem(*UpdateItem) (*UpdateItemResult, error)
	CreateTable(*schema.CreateRequest) (*schema.CreateResponse, error)
}

type defaultExecutor struct {
	endpoint string
	caller   http.Client
	aws      AWSInfo
}

/*
Override the endpoint for this client.

Mostly this is for testing and the defaults should suffice in production.
*/
func (e *defaultExecutor) SetEndpoint(endpoint string) {
	e.endpoint = endpoint
}

func (e *defaultExecutor) makeRequest(target string, document interface{}) ([]byte, error) {
	buf, err := json.Marshal(document)
	// log.Printf("Request Body: \n%s\n\n", buf)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", e.endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-amz-target", DynamoTargetPrefix+target)
	req.Header.Add("content-type", "application/x-amz-json-1.0")
	req.Header.Set("Host", req.URL.Host)
	e.aws.signRequest(req, buf)
	response, err := e.caller.Do(req)
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

func (e *defaultExecutor) makeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error) {
	body, err := e.makeRequest(method, document)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, dest)
	return
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
