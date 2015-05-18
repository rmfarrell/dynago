package dynago

type updateItemRequest struct {
	Key       Document
	TableName string

	ConditionExpression       string   `json:",omitempty"`
	UpdateExpression          string   `json:",omitempty"`
	ExpressionAttributeValues Document `json:",omitempty"`

	ReturnConsumedCapacity      string       `json:",omitempty"` // TODO
	ReturnItemCollectionMetrics string       `json:",omitempty"` // TODO
	ReturnValues                ReturnValues `json:",omitempty"`
}

type updateItemResponse struct {
	Attributes Document
	// TODO ConsumedCapacity
	// TODO ItemCollectionMetrics
}

func newUpdateItem(client *Client, table string, key Document) *UpdateItem {
	return &UpdateItem{
		client: client,
		req: updateItemRequest{
			Key:       key,
			TableName: table,
		},
	}
}

type UpdateItem struct {
	client *Client
	req    updateItemRequest
}

func (u *UpdateItem) ConditionExpression(expression string) *UpdateItem {
	u.req.ConditionExpression = expression
	return u
}

func (u *UpdateItem) UpdateExpression(expression string) *UpdateItem {
	u.req.UpdateExpression = expression
	return u
}

func (u *UpdateItem) Param(key string, value interface{}) *UpdateItem {
	paramHelper(&u.req.ExpressionAttributeValues, key, value)
	return u
}

func (u *UpdateItem) ReturnValues(returnValues ReturnValues) *UpdateItem {
	u.req.ReturnValues = returnValues
	return u
}

func (u *UpdateItem) Execute() (res *UpdateItemResult, err error) {
	if u.req.ReturnValues != ReturnNone && u.req.ReturnValues != "" {
		res = &UpdateItemResult{}
		err = u.client.makeRequestUnmarshal("UpdateItem", &u.req, res)
	} else {
		_, err = u.client.makeRequest("UpdateItem", &u.req)
	}
	return
}

type UpdateItemResult struct {
	Attributes Document
}
