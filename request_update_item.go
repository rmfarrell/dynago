package dynago

type updateItemRequest struct {
	Key       Document
	TableName string

	ConditionExpression string `json:",omitempty"`
	UpdateExpression    string `json:",omitempty"`
	expressionAttributes

	ReturnConsumedCapacity      string       `json:",omitempty"` // TODO
	ReturnItemCollectionMetrics string       `json:",omitempty"` // TODO
	ReturnValues                ReturnValues `json:",omitempty"`
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

// UpdateItem is used to modify a single item in a table.
type UpdateItem struct {
	client *Client
	req    updateItemRequest
}

// ConditionExpression sets a condition which if not satisfied, the PutItem is not performed.
func (u UpdateItem) ConditionExpression(expression string, params ...Params) *UpdateItem {
	u.req.paramsHelper(params)
	u.req.ConditionExpression = expression
	return &u
}

/*
Key allows the key to be changed/set after the UpdateItem has been created.

This allows a generic UpdateItem to be created with all the settings and
expressions, and then re-used for multiple requests with different keys.
*/
func (u UpdateItem) Key(key Document) *UpdateItem {
	u.req.Key = key
	return &u
}

/*
UpdateExpression defines the operations which will be performed by this update.

	UpdateExpression("SET Field1=:val1, Field2=:val2 DELETE Field3")

Expression values cannot be provided inside the expression strings, they must
be referenced by the :param values.

See http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.Modifying.html
for more information on how to use update expression strings
*/
func (u UpdateItem) UpdateExpression(expression string, params ...Params) *UpdateItem {
	u.req.paramsHelper(params)
	u.req.UpdateExpression = expression
	return &u
}

// Param is a shortcut to set a single bound parameter.
func (u UpdateItem) Param(key string, value interface{}) *UpdateItem {
	u.req.paramHelper(key, value)
	return &u
}

// Params sets multiple bound parameters on this query.
func (u UpdateItem) Params(params ...Params) *UpdateItem {
	u.req.paramsHelper(params)
	return &u
}

// ReturnValues allows you to get return values of either updated or old fields (see ReturnValues const)
func (u UpdateItem) ReturnValues(returnValues ReturnValues) *UpdateItem {
	u.req.ReturnValues = returnValues
	return &u
}

// Execute this UpdateItem and return the result.
func (u *UpdateItem) Execute() (res *UpdateItemResult, err error) {
	return u.client.executor.UpdateItem(u)
}

// UpdateItem on this executor.
func (e *AwsExecutor) UpdateItem(u *UpdateItem) (result *UpdateItemResult, err error) {
	if u.req.ReturnValues != ReturnNone && u.req.ReturnValues != "" {
		err = e.MakeRequestUnmarshal("UpdateItem", &u.req, &result)
	} else {
		_, err = e.makeRequest("UpdateItem", &u.req)
	}
	return
}

// UpdateItemResult is returned when ReturnValues is set on the UpdateItem.
type UpdateItemResult struct {
	Attributes Document
}
