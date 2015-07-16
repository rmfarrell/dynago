package dynago

import (
	"github.com/underarmour/dynago/schema"
)

type awsSchemaExecutor struct {
	*AwsExecutor
}

func (e awsSchemaExecutor) CreateTable(req *schema.CreateRequest) (*schema.CreateResult, error) {
	resp := &schema.CreateResult{}
	err := e.MakeRequestUnmarshal("CreateTable", req, resp)
	return resp, err
}

func (e awsSchemaExecutor) DeleteTable(req *schema.DeleteRequest) (*schema.DeleteResult, error) {
	resp := &schema.DeleteResult{}
	err := e.MakeRequestUnmarshal("DeleteTable", req, resp)
	return resp, err
}

func (e awsSchemaExecutor) DescribeTable(req *schema.DescribeRequest) (resp *schema.DescribeResponse, err error) {
	resp = &schema.DescribeResponse{}
	err = e.MakeRequestUnmarshal("DescribeTable", req, resp)
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

func (e awsSchemaExecutor) ListTables(list *ListTables) (*schema.ListResponse, error) {
	resp := &schema.ListResponse{}
	err := e.MakeRequestUnmarshal("ListTables", list.req, resp)
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
