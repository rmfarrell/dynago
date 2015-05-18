package dynago

type queryRequest struct {
	TableName string
	IndexName string `json:",omitempty"`

	// Filtering and query expressions
	KeyConditionExpression    string   `json:",omitempty"`
	FilterExpression          string   `json:",omitempty"`
	ProjectionExpression      string   `json:",omitempty"`
	ExpressionAttributeValues Document `json:",omitempty"`

	CapacityDetail   CapacityDetail `json:"ReturnConsumedCapacity,omitempty"`
	ConsistentRead   *bool          `json:",omitempty"`
	ScanIndexForward *bool          `json:",omitempty"`
	Limit            uint32         `json:",omitempty"`
}

type queryResponse struct {
	//ConsumedCapacity *ConsumedCapacityResponse  // TODO
	Count            int
	Items            []Document
	LastEvaluatedKey *Document
}

func newQuery(client *Client, table string) *Query {
	req := queryRequest{
		TableName: table,
	}
	return &Query{client, req}
}

type Query struct {
	client *Client
	req    queryRequest
}

func (q *Query) ConsistentRead(strong bool) *Query {
	q.req.ConsistentRead = &strong
	return q
}

func (q *Query) FilterExpression(expression string) *Query {
	q.req.FilterExpression = expression
	return q
}

func (q *Query) KeyConditionExpression(expression string) *Query {
	q.req.KeyConditionExpression = expression
	return q
}

func (q *Query) ProjectionExpression(expression string) *Query {
	q.req.ProjectionExpression = expression
	return q
}

// Shortcut to set a single parameter for ExpressionAttributeValues.
func (q *Query) Param(key string, value interface{}) *Query {
	paramHelper(&q.req.ExpressionAttributeValues, key, value)
	return q
}

func (q *Query) Desc() *Query {
	forward := false
	q.req.ScanIndexForward = &forward
	return q
}

func (q *Query) Execute() (result *QueryResult, err error) {
	var response queryResponse
	err = q.client.makeRequestUnmarshal("Query", &q.req, &response)
	if err != nil {
		return
	}
	result = &QueryResult{
		Items: response.Items,
		Count: response.Count,
	}
	return
}

type QueryResult struct {
	Items []Document
	Count int
}

// Helper for a variety of endpoint types to build a params dictionary.
func paramHelper(doc *Document, key string, value interface{}) {
	if *doc == nil {
		*doc = Document{}
	}
	(*doc)[key] = value
}
