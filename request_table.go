package dynago

import (
	"github.com/underarmour/dynago/schema"
)

func (e *awsExecutor) CreateTable(req *schema.CreateRequest) (*schema.CreateResponse, error) {
	resp := &schema.CreateResponse{}
	err := e.makeRequestUnmarshal("CreateTable", req, resp)
	return resp, err
}
