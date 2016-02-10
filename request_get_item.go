package dynago

type getItemRequest struct {
	TableName string

	Key Document

	ProjectionExpression string `json:",omitempty"`
	expressionAttributes

	ReturnConsumedCapacity CapacityDetail `json:",omitempty"`
	ConsistentRead         bool           `json:",omitempty"`
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

// GetItem is used to get a single item by key from the table.
type GetItem struct {
	client *Client
	req    getItemRequest
}

// ProjectionExpression allows the client to specify which attributes are returned.
func (p GetItem) ProjectionExpression(expression string, params ...Params) *GetItem {
	p.req.ProjectionExpression = expression
	p.req.paramsHelper(params)
	return &p
}

// Param is a shortcut to set a single bound parameter.
func (p GetItem) Param(key string, value interface{}) *GetItem {
	p.req.paramHelper(key, value)
	return &p
}

// Params sets multiple bound parameters on this query.
func (p GetItem) Params(params ...Params) *GetItem {
	p.req.paramsHelper(params)
	return &p
}

// ReturnConsumedCapacity enables capacity reporting on this GetItem.
// Defaults to CapacityNone if not set
func (p GetItem) ReturnConsumedCapacity(consumedCapacity CapacityDetail) *GetItem {
	p.req.ReturnConsumedCapacity = consumedCapacity
	return &p
}

// ConsistentRead enables strongly consistent reads if the argument is true.
func (p GetItem) ConsistentRead(strong bool) *GetItem {
	p.req.ConsistentRead = strong
	return &p
}

// Execute the get item.
func (p *GetItem) Execute() (result *GetItemResult, err error) {
	return p.client.executor.GetItem(p)
}

// GetItem gets a single item.
func (e *AwsExecutor) GetItem(g *GetItem) (result *GetItemResult, err error) {
	err = e.MakeRequestUnmarshal("GetItem", &g.req, &result)
	return
}

// GetItemResult is the result from executing a GetItem.
type GetItemResult struct {
	Item             Document
	ConsumedCapacity *ConsumedCapacity
}
