package dynago

import (
	"github.com/underarmour/dynago/schema"
)

func (e *awsExecutor) CreateTable(req *schema.CreateRequest) (*schema.CreateResult, error) {
	resp := &schema.CreateResult{}
	err := e.makeRequestUnmarshal("CreateTable", req, resp)
	return resp, err
}

func (e *awsExecutor) DeleteTable(req *schema.DeleteRequest) (*schema.DeleteResult, error) {
	resp := &schema.DeleteResult{}
	err := e.makeRequestUnmarshal("DeleteTable", req, resp)
	return resp, err
}

func (e *awsExecutor) DescribeTable(req *schema.DescribeRequest) (resp *schema.DescribeResponse, err error) {
	resp = &schema.DescribeResponse{}
	err = e.makeRequestUnmarshal("DescribeTable", req, resp)
	return
}

type ListTables struct {
	client *Client
	req    schema.ListRequest
}

func (l ListTables) Limit(limit uint) *ListTables {
	l.req.Limit = limit
	return &l
}

func (l *ListTables) Execute() (result *ListTablesResult, err error) {
	resp, err := l.client.schemaExecutor.ListTables(l)
	if err == nil {
		result = &ListTablesResult{resp.TableNames, resp.LastEvaluatedTableName, l}
	}
	return
}

func (e *awsExecutor) ListTables(list *ListTables) (*schema.ListResponse, error) {
	resp := &schema.ListResponse{}
	err := e.makeRequestUnmarshal("ListTables", list.req, resp)
	return resp, err
}

type ListTablesResult struct {
	TableNames []string
	cursor     *string
	req        *ListTables
}

// Helper to get the ListTables for the next page of listings.
// If there is not a next page, returns nil
func (r ListTablesResult) Next() *ListTables {
	if r.cursor == nil {
		return nil
	}
	return &ListTables{r.req.client, schema.ListRequest{
		ExclusiveStartTableName: *r.cursor,
		Limit: r.req.req.Limit,
	}}
}
