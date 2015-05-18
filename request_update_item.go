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

// Set a condition expression for conditional update.
func (u *UpdateItem) ConditionExpression(expression string) *UpdateItem {
	u.req.ConditionExpression = expression
	return u
}

// Set an update expression to update specific fields and values.
func (u *UpdateItem) UpdateExpression(expression string) *UpdateItem {
	u.req.UpdateExpression = expression
	return u
}

// Quick-set parameters for all the values here.
func (u *UpdateItem) Param(key string, value interface{}) *UpdateItem {
	paramHelper(&u.req.ExpressionAttributeValues, key, value)
	return u
}

// If set, then we will get return values of either updated or old fields (see ReturnValues const)
func (u *UpdateItem) ReturnValues(returnValues ReturnValues) *UpdateItem {
	u.req.ReturnValues = returnValues
	return u
}

// Execute this UpdateItem and return the result.
func (u *UpdateItem) Execute() (res *UpdateItemResult, err error) {
	return u.client.executor.UpdateItem(u)
}

func (e *awsExecutor) UpdateItem(u *UpdateItem) (res *UpdateItemResult, err error) {
	if u.req.ReturnValues != ReturnNone && u.req.ReturnValues != "" {
		res = &UpdateItemResult{}
		err = e.makeRequestUnmarshal("UpdateItem", &u.req, res)
	} else {
		_, err = e.makeRequest("UpdateItem", &u.req)
	}
	return
}

type UpdateItemResult struct {
	Attributes Document
}
