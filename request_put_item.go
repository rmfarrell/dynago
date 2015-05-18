package dynago

type putItemRequest struct {
	TableName string

	Item Document

	ConditionExpression       string   `json:",omitempty"`
	ExpressionAttributeValues Document `json:",omitempty"`

	// TODO ReturnConsumedCapacity string
	// TODO ReturnItemCollectionMetrics
	ReturnValues ReturnValues `json:",omitempty"`
}

func newPutItem(client *Client, table string, item Document) *PutItem {
	return &PutItem{
		client: client,
		req: putItemRequest{
			TableName: table,
			Item:      item,
		},
	}
}

type PutItem struct {
	client *Client
	req    putItemRequest
}

func (p *PutItem) ConditionExpression(expression string) *PutItem {
	p.req.ConditionExpression = expression
	return p
}

func (p *PutItem) Param(key string, value interface{}) *PutItem {
	paramHelper(&p.req.ExpressionAttributeValues, key, value)
	return p
}

func (p *PutItem) ReturnValues(returnValues ReturnValues) *PutItem {
	p.req.ReturnValues = returnValues
	return p
}

/*
Actually execute this put.

PutItemResult may be empty if there is no need
*/
func (p *PutItem) Execute() (res *PutItemResult, err error) {
	res = &PutItemResult{}
	if p.req.ReturnValues != ReturnNone && p.req.ReturnValues != "" {
		err = p.client.makeRequestUnmarshal("PutItem", &p.req, res)
	} else {
		_, err = p.client.makeRequest("PutItem", &p.req)
	}
	return
}

type PutItemResult struct {
	Attributes Document
}
