package streams

const targetPrefix = "DynamoDBStreams_20120810." // This is the Dynamo API version we support

// Config is configuration for a streams client.
type Config struct {
	// Must implement MakeRequestUnmarshal. The easiest thing to use in this
	// case is an instance of dynago.AwsExecutor or the like.
	Requester MakeRequester
}

func NewClient(config *Config) *Client {
	return &Client{config.Requester}
}

/*
Client is the low-level interface to the streams API.

Here, all the individual Streams API requests are available.
Because this is a low-level client package, the main aim is to provide mappings
for the core API calls needed to implement streaming.

For more information on how the API's interact, check out:
http://docs.aws.amazon.com/dynamodbstreams/latest/APIReference/Welcome.html
*/
type Client struct {
	caller MakeRequester
}

// DescribeStream gets metadata about a stream and information about its shards.
func (s *Client) DescribeStream(request *DescribeStreamRequest) (dest *DescribeStreamResponse, err error) {
	err = s.caller.MakeRequestUnmarshal(targetPrefix+"DescribeStream", request, &dest)
	return
}

// GetShardIterator is required to get a valid iterator before you begin streaming.
func (s *Client) GetShardIterator(req *GetIteratorRequest) (dest *GetIteratorResult, err error) {
	err = s.caller.MakeRequestUnmarshal(targetPrefix+"GetShardIterator", req, &dest)
	return
}

// GetRecords will get the next set of records inside a DynamoDB streams shard.
//
// An empty set of records does not mean the stream is complete; only the
// absence of an iterator signals completion of a stream shard (such as the
// shard is now closed for new data)
func (s *Client) GetRecords(req *GetRecordsRequest) (result *GetRecordsResponse, err error) {
	err = s.caller.MakeRequestUnmarshal(targetPrefix+"GetRecords", req, &result)
	return
}

// MakeRequester is Equivalent to the dynago version
type MakeRequester interface {
	MakeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error)
}
