package dynago

/*
ConsumedCapacity describes the write or read capacity units consumed by a
given operation.

This is provided in a variety DynamoDB responses when consumed capacity
detail is explicitly requested, and various elements may be unset or be nil
based on the response.

Not all operations set all fields, so expect some of the fields to be empty
or nil.
*/
type ConsumedCapacity struct {
	// Total capacity units for the entire operation
	CapacityUnits float64

	// Capacity units from local or global secondary indexes
	GlobalSecondaryIndexes map[string]*Capacity
	LocalSecondaryIndexes  map[string]*Capacity

	// Capacity units by table
	TableName string
	Table     *Capacity
}

// Capacity descriptor for indexes and tables.
type Capacity struct {
	CapacityUnits float64
}

// BatchConsumedCapacity describes ConsumedCapacity for multiple tables in batch requests.
type BatchConsumedCapacity []ConsumedCapacity

/*
GetTable returns the batch consumed capacity for the named table.
*/
func (b BatchConsumedCapacity) GetTable(tableName string) *ConsumedCapacity {
	for i := range b {
		c := &b[i]
		if c.TableName == tableName {
			return c
		}
	}
	return nil
}
