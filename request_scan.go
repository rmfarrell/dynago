package dynago

type scanRequest struct {
	queryRequest
	Segment       *int `json:",omitempty"`
	TotalSegments *int `json:",omitempty"`
}

type Scan struct {
	client *Client
	req    scanRequest
}

func newScan(c *Client, table string) *Scan {
	req := scanRequest{queryRequest: queryRequest{TableName: table}}
	return &Scan{c, req}
}

/*
ExclusiveStartKey sets the start key (effectively the offset cursor).
*/
func (s Scan) ExclusiveStartKey(key Document) *Scan {
	s.req.ExclusiveStartKey = key
	return &s
}

/*
FilterExpression post-filters results on this scan.

Scans with a FilterExpression may return 0 results due to scanning past
records which don't match the filter, but still have more results to get.
*/
func (s Scan) FilterExpression(expression string, params ...Params) *Scan {
	s.req.FilterExpression = expression
	s.req.paramsHelper(params)
	return &s
}

// IndexName specifies a secondary index to scan instead of a table.
func (s Scan) IndexName(name string) *Scan {
	s.req.IndexName = name
	return &s
}

// Limit the maximum number of results to return per call.
func (s Scan) Limit(limit uint) *Scan {
	s.req.Limit = limit
	return &s
}

// ProjectionExpression allows the client to specify which attributes are returned.
func (s Scan) ProjectionExpression(expression string, params ...Params) *Scan {
	s.req.ProjectionExpression = expression
	s.req.paramsHelper(params)
	return &s
}

// ReturnConsumedCapacity enables capacity reporting on this Query.
func (s Scan) ReturnConsumedCapacity(consumedCapacity CapacityDetail) *Scan {
	s.req.ReturnConsumedCapacity = consumedCapacity
	return &s
}

// Segment chooses the parallel segment of the table to scan.
func (s Scan) Segment(segment, total int) *Scan {
	s.req.Segment = &segment
	s.req.TotalSegments = &total
	return &s
}

/*
Select specifies how attributes are chosen, or enables count mode.

Most of the time, specifying Select is not required, because the DynamoDB
API does the "right thing" inferring values based on other attributes like
the projection expression, index, etc.
*/
func (s Scan) Select(value Select) *Scan {
	s.req.Select = value
	return &s
}

// Execute this Scan query.
func (s *Scan) Execute() (*ScanResult, error) {
	return s.client.executor.Scan(s)
}

// Scan operation
func (e *AwsExecutor) Scan(s *Scan) (result *ScanResult, err error) {
	result = &ScanResult{req: s}
	err = e.MakeRequestUnmarshal("Scan", s.req, &result)
	return
}

// ScanResult is the result of Scan queries.
type ScanResult struct {
	req              *Scan
	Items            []Document
	LastEvaluatedKey Document
	ConsumedCapacity *ConsumedCapacity
}

/*
Next returns a scan which get the next page of results when executed.

If the scan has a LastEvaluatedKey, returns another Scan. Otherwise, returns nil.
*/
func (r *ScanResult) Next() *Scan {
	if r.LastEvaluatedKey != nil {
		return r.req.ExclusiveStartKey(r.LastEvaluatedKey)
	}
	return nil
}
