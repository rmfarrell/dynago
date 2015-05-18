package dynago

type getItemRequest struct {
	TableName string

	Key Document

	ProjectionExpression      string   `json:",omitempty"`
	ExpressionAttributeValues Document `json:",omitempty"`

	// TODO ReturnConsumedCapacity string
	ConsistentRead bool `json:",omitempty"`
}

func newGetItem(client *Client, table string, key Document) *GetItem {
	return &GetItem{
		client: client,
		req: getItemRequest{
			TableName: table,
			Key:       key,
		},
	}
}

type GetItem struct {
	client *Client
	req    getItemRequest
}

func (p *GetItem) ProjectionExpression(expression string) *GetItem {
	p.req.ProjectionExpression = expression
	return p
}

func (p *GetItem) Param(key string, value interface{}) *GetItem {
	paramHelper(&p.req.ExpressionAttributeValues, key, value)
	return p
}

func (p *GetItem) ConsistentRead() *GetItem {
	p.req.ConsistentRead = true
	return p
}

func (p *GetItem) Execute() (result *GetItemResult, err error) {
	result = &GetItemResult{}
	err = p.client.makeRequestUnmarshal("GetItem", &p.req, result)
	return
}

type GetItemResult struct {
	Item Document
}
