package dynago

type putItemRequest struct {
	TableName string
	Item      Document

	ConditionExpression string `json:",omitempty"`
	expressionAttributes

	// TODO ReturnItemCollectionMetrics
	ReturnConsumedCapacity CapacityDetail `json:",omitempty"`
	ReturnValues           ReturnValues   `json:",omitempty"`
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

// PutItem is used to create/replace single items in the table.
type PutItem struct {
	client *Client
	req    putItemRequest
}

// ConditionExpression sets a condition which if not satisfied, the PutItem is not performed.
func (p PutItem) ConditionExpression(expression string, params ...Params) *PutItem {
	p.req.ConditionExpression = expression
	p.req.paramsHelper(params)
	return &p
}

// Param is a shortcut to set a single bound parameter.
func (p PutItem) Param(key string, value interface{}) *PutItem {
	p.req.paramHelper(key, value)
	return &p
}

// Params sets multiple bound parameters on this query.
func (p PutItem) Params(params ...Params) *PutItem {
	p.req.paramsHelper(params)
	return &p
}

// ReturnConsumedCapacity enables capacity reporting on this PutItem.
func (p PutItem) ReturnConsumedCapacity(consumedCapacity CapacityDetail) *PutItem {
	p.req.ReturnConsumedCapacity = consumedCapacity
	return &p
}

// ReturnValues can allow you to ask for either previous or new values on an update
func (p PutItem) ReturnValues(returnValues ReturnValues) *PutItem {
	p.req.ReturnValues = returnValues
	return &p
}

/*
Execute this PutItem.

PutItemResult will be nil unless ReturnValues or ReturnConsumedCapacity is set.
*/
func (p *PutItem) Execute() (res *PutItemResult, err error) {
	return p.client.executor.PutItem(p)
}

// PutItem on this executor.
func (e *AwsExecutor) PutItem(p *PutItem) (res *PutItemResult, err error) {
	if (p.req.ReturnValues != ReturnNone && p.req.ReturnValues != "") || p.req.ReturnConsumedCapacity != "" {
		err = e.MakeRequestUnmarshal("PutItem", &p.req, &res)
	} else {
		_, err = e.makeRequest("PutItem", &p.req)
	}
	return
}

// PutItemResult is returned when a PutItem is executed.
type PutItemResult struct {
	Attributes       Document
	ConsumedCapacity *ConsumedCapacity
}
