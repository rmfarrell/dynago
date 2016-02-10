package schema

// CreateRequest is the request to create a new DynamoDB table.
type CreateRequest struct {
	TableName              string
	AttributeDefinitions   []AttributeDefinition
	KeySchema              []KeySchema
	ProvisionedThroughput  ProvisionedThroughput
	GlobalSecondaryIndexes []SecondaryIndex
	LocalSecondaryIndexes  []SecondaryIndex
	StreamSpecification    *StreamSpecification `json:",omitempty"`
}

func NewCreateRequest(table string) *CreateRequest {
	return &CreateRequest{
		TableName:             table,
		ProvisionedThroughput: ProvisionedThroughput{1, 1},
	}
}

// Simple way to add a hash key and attribute definition in one go.
func (r *CreateRequest) HashKey(name string, attributeType AttributeType) *CreateRequest {
	r.ensureAttribute(name, attributeType)
	r.KeySchema = append(r.KeySchema, KeySchema{name, HashKey})
	return r
}

// RangeKey is a simple way to add a hash key and attribute definition in one go.
func (r *CreateRequest) RangeKey(name string, attributeType AttributeType) *CreateRequest {
	r.ensureAttribute(name, attributeType)
	r.KeySchema = append(r.KeySchema, KeySchema{name, RangeKey})
	return r
}

func (r *CreateRequest) ensureAttribute(name string, attributeType AttributeType) {
	for _, a := range r.AttributeDefinitions {
		if a.AttributeName == name {
			return
		}
	}
	r.AttributeDefinitions = append(r.AttributeDefinitions, AttributeDefinition{name, attributeType})
}

// CreateResult describes the table created
type CreateResult struct {
	TableDescription TableDescription
}

// DeleteRequest asks to delete a DynamoDB table
type DeleteRequest struct {
	TableName string
}

// DeleteResult describes the table deleted
type DeleteResult struct {
	TableDescription TableDescription
}

// ListRequest asks to list all tables on the current account.
type ListRequest struct {
	ExclusiveStartTableName string `json:",omitempty"`
	Limit                   uint   `json:",omitempty"`
}

// ListResponse is the response of table list request.
type ListResponse struct {
	LastEvaluatedTableName *string
	TableNames             []string
}

// DescribeRequest gives details about a single table.
type DescribeRequest struct {
	TableName string
}

type DescribeResponse struct {
	Table TableDescription
}
