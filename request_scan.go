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

func (s Scan) ExclusiveStartKey(key Document) *Scan {
	s.req.ExclusiveStartKey = key
	return &s
}

func (s Scan) FilterExpression(expression string, params ...Params) *Scan {
	s.req.FilterExpression = expression
	s.req.paramsHelper(params)
	return &s
}

func (s Scan) IndexName(name string) *Scan {
	s.req.IndexName = name
	return &s
}

func (s Scan) Limit(limit uint) *Scan {
	s.req.Limit = limit
	return &s
}

func (s Scan) ProjectionExpression(expression string, params ...Params) *Scan {
	s.req.ProjectionExpression = expression
	s.req.paramsHelper(params)
	return &s
}

// Choose the parallel segment of the table to scan.
func (s Scan) Segment(segment, total int) *Scan {
	s.req.Segment = &segment
	s.req.TotalSegments = &total
	return &s
}

func (s *Scan) Execute() (*ScanResult, error) {
	return s.client.executor.Scan(s)
}

func (e *awsExecutor) Scan(s *Scan) (result *ScanResult, err error) {
	result = &ScanResult{req: s}
	err = e.makeRequestUnmarshal("Scan", s.req, &result)
	return
}

type ScanResult struct {
	req              *Scan
	Items            []Document
	LastEvaluatedKey Document
}

func (r *ScanResult) Next() *Scan {
	if r.LastEvaluatedKey != nil {
		return r.req.ExclusiveStartKey(r.LastEvaluatedKey)
	}
	return nil
}
