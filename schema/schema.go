package schema

type KeyType string

const (
	HashKey  KeyType = "HASH"
	RangeKey KeyType = "RANGE"
)

type AttributeType string

const (
	String AttributeType = "S"
	Number AttributeType = "N"
	Binary AttributeType = "B"
)

type TableDescription struct {
	TableName        string
	TableSizeBytes   uint64
	TableStatus      string
	CreationDateTime float64

	AttributeDefinitions   []AttributeDefinition
	GlobalSecondaryIndexes []SecondaryIndexResponse
}

type ProvisionedThroughput struct {
	ReadCapacityUnits  uint
	WriteCapacityUnits uint
}

type AttributeDefinition struct {
	AttributeName string
	AttributeType AttributeType
}

type KeySchema struct {
	AttributeName string
	KeyType       KeyType
}

type SecondaryIndex struct {
	IndexName             string
	KeySchema             []KeySchema
	Projection            interface{} // TODO
	ProvisionedThroughput ProvisionedThroughput
}

// Secondary indexes as described in table descriptions
type SecondaryIndexResponse struct {
	SecondaryIndex
	Backfilling    bool
	IndexSizeBytes int
	IndexStatus    string
	ItemCount      int
}
