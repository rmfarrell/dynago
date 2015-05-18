package dynatest

import (
	"github.com/crast/dynago"
	"github.com/crast/dynago/schema"
)

// A Mock executor
type Executor struct {
}

func (e *Executor) GetItem(*dynago.GetItem) (*dynago.GetItemResult, error) {
	return nil, nil
}

func (e *Executor) PutItem(*dynago.PutItem) (*dynago.PutItemResult, error) {
	return nil, nil
}

func (e *Executor) Query(*dynago.Query) (*dynago.QueryResult, error) {
	return nil, nil
}

func (e *Executor) UpdateItem(*dynago.UpdateItem) (*dynago.UpdateItemResult, error) {
	return nil, nil
}

func (e *Executor) CreateTable(*schema.CreateRequest) (*schema.CreateResponse, error) {
	return nil, nil
}
