package streams

import (
	"github.com/rmfarrell/dynago"
	"github.com/rmfarrell/dynago/schema"
)

// Stream is the compact representation of a stream.
type Stream struct {
	StreamArn   string
	StreamLabel string
	TableName   string
}

// StreamDescription is the main response value from DescribeStream
type StreamDescription struct {
	Stream
	KeySchema      []schema.KeySchema
	Shards         []Shard
	StreamStatus   string
	StreamViewType string

	CreationRequestDateTime float64
	LastEvaluatedShardId    string
}

// Shard describes one of the shards of a stream
type Shard struct {
	ParentShardId       string
	SequenceNumberRange SequenceNumberRange
	ShardId             string
}

// SequenceNumberRange is information about sequence numbers in a stream
type SequenceNumberRange struct {
	EndingSequenceNumber   string
	StartingSequenceNumber string
}

// StreamRecord  is a description of a single data modification that was performed on an item in a DynamoDB table.
type StreamRecord struct {
	Keys           dynago.Document
	OldImage       dynago.Document
	NewImage       dynago.Document
	SequenceNumber string
	SizeBytes      uint64
	StreamViewType string
}

// Record is a description of a unique event within a stream.
type Record struct {
	StreamRecord `json:"dynamodb"`
	AwsRegion    string
	EventId      string
	EventName    string
	EventSource  string
	EventVersion string
}
