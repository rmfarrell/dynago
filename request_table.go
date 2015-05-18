package dynago

import (
	"github.com/crast/dynago/schema"
)

func (e *awsExecutor) CreateTable(req *schema.CreateRequest) (*schema.CreateResponse, error) {
	resp := &schema.CreateResponse{}
	err := e.makeRequestUnmarshal("CreateTable", req, resp)
	return resp, err
}
