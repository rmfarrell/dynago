package dynago

type batchWriteItemRequest struct {
	RequestItems BatchWriteTableMap

	ReturnConsumedCapacity CapacityDetail `json:",omitempty"`
}

// BatchWriteTableMap describes writes where the key is the table and the value is the unprocessed items.
type BatchWriteTableMap map[string][]*BatchWriteTableEntry

// BatchWriteTableEntry is a single write or delete request.
type BatchWriteTableEntry struct {
	DeleteRequest *batchDelete `json:",omitempty"`
	PutRequest    *batchPut    `json:",omitempty"`
}

// SetDelete sets this table entry as a delete request
func (e *BatchWriteTableEntry) SetDelete(key Document) {
	e.DeleteRequest = &batchDelete{key}
}

// SetPut sets this table entry as a put request
func (e *BatchWriteTableEntry) SetPut(item Document) {
	e.PutRequest = &batchPut{item}
}

type batchDelete struct {
	Key Document
}

type batchPut struct {
	Item Document
}

type batchAction struct {
	next  *batchAction
	table string
	item  Document
}

func newBatchWrite(client *Client) *BatchWrite {
	return &BatchWrite{
		client: client,
	}
}

// BatchWrite allows writing many items in a single roundtrip to DynamoDB.
type BatchWrite struct {
	client  *Client
	puts    *batchAction
	deletes *batchAction

	capacityDetail CapacityDetail
}

/*
Put queues up some number of puts to a table.
*/
func (b BatchWrite) Put(table string, items ...Document) *BatchWrite {
	addBatchActions(&b.puts, table, items)
	return &b
}

/*
Delete queues some number of deletes for a table.
*/
func (b BatchWrite) Delete(table string, keys ...Document) *BatchWrite {
	addBatchActions(&b.deletes, table, keys)
	return &b
}

// ReturnConsumedCapacity enables capacity reporting on this Query.
func (b BatchWrite) ReturnConsumedCapacity(consumedCapacity CapacityDetail) *BatchWrite {
	b.capacityDetail = consumedCapacity
	return &b
}

// Execute the writes in this batch.
func (b *BatchWrite) Execute() (*BatchWriteResult, error) {
	return b.client.executor.BatchWriteItem(b)
}

// Build the table map that is represented by this BatchWrite
func (b *BatchWrite) buildTableMap() (m BatchWriteTableMap) {
	m = BatchWriteTableMap{}
	ensure := func(table string) (r *BatchWriteTableEntry) {
		r = &BatchWriteTableEntry{}
		m[table] = append(m[table], r)
		return
	}

	for put := b.puts; put != nil; put = put.next {
		ensure(put.table).SetPut(put.item)
	}

	for d := b.deletes; d != nil; d = d.next {
		ensure(d.table).SetDelete(d.item)
	}
	return
}

// BatchWriteItem executes multiple puts/deletes in a single roundtrip.
func (e *AwsExecutor) BatchWriteItem(b *BatchWrite) (result *BatchWriteResult, err error) {
	req := batchWriteItemRequest{
		RequestItems:           b.buildTableMap(),
		ReturnConsumedCapacity: b.capacityDetail,
	}

	err = e.MakeRequestUnmarshal("BatchWriteItem", req, &result)
	return
}

// BatchWriteResult explains what happened in a batch write.
type BatchWriteResult struct {
	UnprocessedItems BatchWriteTableMap
	ConsumedCapacity BatchConsumedCapacity
}

///////////////////// Batch Get

const (
	bgProjectionExpression = "ProjectionExpression"
	bgProjectionParams     = "Params"
	bgConsistentRead       = "ConsistentRead"
)

type batchGetItemRequest struct {
	RequestItems BatchGetTableMap

	ReturnConsumedCapacity CapacityDetail `json:",omitempty"`
}

type BatchGetTableMap map[string]*BatchGetTableEntry

type BatchGetTableEntry struct {
	Keys []Document

	expressionAttributes
	ProjectionExpression string `json:",omitempty"`
	ConsistentRead       bool   `json:",omitempty"`
}

// BatchGet allows getting multiple items by key from a table.
type BatchGet struct {
	client  *Client
	gets    *batchAction
	options *batchAction

	capacityDetail CapacityDetail
}

/*
Get queues some gets for a table.
Can be called multiple times to queue up gets for multiple tables.
*/
func (b BatchGet) Get(table string, keys ...Document) *BatchGet {
	addBatchActions(&b.gets, table, keys)
	return &b
}

/*
ProjectionExpression allows the client to specify attributes returned for a table.

Projection expression is scoped to each table, and must be called for each
table on which you want a ProjectionExpression.
*/
func (b BatchGet) ProjectionExpression(table string, expression string, params ...Params) *BatchGet {
	doc := Document{
		bgProjectionExpression: expression,
		bgProjectionParams:     params,
	}
	addBatchActions(&b.options, table, []Document{doc})
	return &b
}

/*
ConsistentRead enables strongly consistent reads per-table.

Consistent read is scoped to each table, so must be called for each table in
this BatchGet for which you want consistent reads.
*/
func (b BatchGet) ConsistentRead(table string, consistent bool) *BatchGet {
	doc := Document{
		bgConsistentRead: consistent,
	}
	addBatchActions(&b.options, table, []Document{doc})
	return &b
}

// ReturnConsumedCapacity enables capacity reporting on this BatchGet
func (b BatchGet) ReturnConsumedCapacity(consumedCapacity CapacityDetail) *BatchGet {
	b.capacityDetail = consumedCapacity
	return &b
}

func (b *BatchGet) buildTableMap() BatchGetTableMap {
	m := BatchGetTableMap{}
	ensure := func(key string) (entry *BatchGetTableEntry) {
		if entry = m[key]; entry == nil {
			entry = &BatchGetTableEntry{}
			m[key] = entry
		}
		return
	}
	for get := b.gets; get != nil; get = get.next {
		entry := ensure(get.table)
		entry.Keys = append(entry.Keys, get.item)
	}
	for option := b.options; option != nil; option = option.next {
		entry := ensure(option.table)
		for k, v := range option.item {
			switch k {
			case bgProjectionExpression:
				entry.ProjectionExpression = v.(string)
			case bgProjectionParams:
				entry.paramsHelper(v.([]Params))
			case bgConsistentRead:
				entry.ConsistentRead = v.(bool)
			}
		}
	}
	return m
}

// Execute this batch get.
func (b *BatchGet) Execute() (result *BatchGetResult, err error) {
	return b.client.executor.BatchGetItem(b)
}

// BatchGetItem gets multiple keys.
func (e *AwsExecutor) BatchGetItem(b *BatchGet) (result *BatchGetResult, err error) {
	req := batchGetItemRequest{
		RequestItems:           b.buildTableMap(),
		ReturnConsumedCapacity: b.capacityDetail,
	}
	err = e.MakeRequestUnmarshal("BatchGetItem", &req, &result)
	return
}

// BatchGetResult is the result of a batch get.
type BatchGetResult struct {
	// Responses to the Batch Get query which were resolved.
	// Note that the order of documents is not guaranteed to be the same as
	// the order of documents requested in building this query.
	Responses map[string][]Document // table name -> list of items.

	// Unprocessed keys are keys which for some reason could not be retrieved.
	// This could be because of response size exceeding the limit or the
	// provisioned throughput being exceeded on one or more tables in this request.
	UnprocessedKeys BatchGetTableMap // Table name -> keys and settings

	ConsumedCapacity BatchConsumedCapacity
}

func addBatchActions(list **batchAction, table string, items []Document) {
	head := *list
	for _, item := range items {
		head = &batchAction{head, table, item}
	}
	*list = head
}
