package dynago

type queryRequest struct {
	TableName string
	IndexName string `json:",omitempty"`

	// Filtering and query expressions
	KeyConditionExpression string `json:",omitempty"`
	FilterExpression       string `json:",omitempty"`
	ProjectionExpression   string `json:",omitempty"`
	expressionAttributes

	CapacityDetail   CapacityDetail `json:"ReturnConsumedCapacity,omitempty"`
	ConsistentRead   *bool          `json:",omitempty"`
	ScanIndexForward *bool          `json:",omitempty"`

	// Limit/offset
	Limit             uint     `json:",omitempty"`
	ExclusiveStartKey Document `json:",omitempty"`
}

type queryResponse struct {
	//ConsumedCapacity *ConsumedCapacityResponse  // TODO
	Count            int
	Items            []Document
	LastEvaluatedKey Document
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

func (q Query) IndexName(name string) *Query {
	q.req.IndexName = name
	return &q
}

// If strong is true, do a strongly consistent read. (defaults to false)
func (q Query) ConsistentRead(strong bool) *Query {
	q.req.ConsistentRead = &strong
	return &q
}

// Set a post-filter expression for the results we scan.
func (q Query) FilterExpression(expression string, params ...Params) *Query {
	q.req.paramsHelper(params)
	q.req.FilterExpression = expression
	return &q
}

// Set a condition expression on the key to narrow down what we scan
func (q Query) KeyConditionExpression(expression string, params ...Params) *Query {
	q.req.paramsHelper(params)
	q.req.KeyConditionExpression = expression
	return &q
}

// Set a Projection Expression for controlling which attributes are returned.
func (q Query) ProjectionExpression(expression string, params ...Params) *Query {
	q.req.paramsHelper(params)
	q.req.ProjectionExpression = expression
	return &q
}

// Shortcut to set a single parameter for ExpressionAttributeValues.
func (q Query) Param(key string, value interface{}) *Query {
	q.req.paramHelper(key, value)
	return &q
}

// Set a param, a document of params, or multiple params
func (q Query) Params(params ...Params) *Query {
	q.req.paramsHelper(params)
	return &q
}

// Return results descending.
func (q Query) Desc() *Query {
	forward := false
	q.req.ScanIndexForward = &forward
	return &q
}

/*
Set the limit on results count.

Note that getting less than `limit` records doesn't mean you're at the last
page of results, that can only be safely asserted if there is no
LastEvaluatedKey on the result.
*/
func (q Query) Limit(limit uint) *Query {
	q.req.Limit = limit
	return &q
}

/*
Set the start key (effectively the offset cursor)
*/
func (q Query) ExclusiveStartKey(key Document) *Query {
	q.req.ExclusiveStartKey = key
	return &q
}

// Execute this query and return results.
func (q *Query) Execute() (result *QueryResult, err error) {
	return q.client.executor.Query(q)
}

func (e *awsExecutor) Query(q *Query) (result *QueryResult, err error) {
	var response queryResponse
	err = e.makeRequestUnmarshal("Query", &q.req, &response)
	if err != nil {
		return
	}
	result = &QueryResult{
		query:            q,
		Items:            response.Items,
		Count:            response.Count,
		LastEvaluatedKey: response.LastEvaluatedKey,
	}
	return
}

// The result returned from a query.
type QueryResult struct {
	query            *Query
	Items            []Document
	Count            int      // The total number of items (for pagination)
	LastEvaluatedKey Document // The offset key for the next page.
}

// Helper for getting a query which will get the next page of results.
// Returns nil if there's no next page.
func (qr *QueryResult) Next() (query *Query) {
	if qr.LastEvaluatedKey != nil && len(qr.LastEvaluatedKey) > 0 {
		query = qr.query.ExclusiveStartKey(qr.LastEvaluatedKey)
	}
	return
}
