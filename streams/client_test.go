package streams_test

import (
	"testing"

	"gopkg.in/underarmour/dynago.v1"
	"gopkg.in/underarmour/dynago.v1/schema"
	"gopkg.in/underarmour/dynago.v1/streams"

	"github.com/stretchr/testify/assert"
)

type mockRequester struct {
	target      string
	body        []byte
	returnBody  []byte
	returnError error
}

func (fr *mockRequester) MakeRequest(target string, body []byte) ([]byte, error) {
	fr.target = target
	fr.body = body
	return fr.returnBody, fr.returnError
}

func setUp(t *testing.T) (*assert.Assertions, *mockRequester, *streams.Client) {
	requester := &mockRequester{}
	executor := &dynago.AwsExecutor{Requester: requester}
	client := streams.NewClient(&streams.Config{
		Requester: executor,
	})
	return assert.New(t), requester, client
}

func TestDescribeStreamEncode(t *testing.T) {
	assert, mock, client := setUp(t)
	mock.returnBody = []byte(`{"StreamDescription":` + streamDescription + `}`)
	result, err := client.DescribeStream(&streams.DescribeStreamRequest{
		StreamArn: "Foo:bar:baz",
		Limit:     50,
	})
	assert.Equal("DynamoDBStreams_20120810.DescribeStream", mock.target)
	assert.Equal(`{"Limit":50,"StreamArn":"Foo:bar:baz"}`, string(mock.body))
	assert.NoError(err)
	sd := result.StreamDescription
	assert.Equal("foo", sd.StreamArn)
	assert.Equal("My Stream", sd.StreamLabel)
	assert.Equal("bar", sd.TableName)
	assert.Equal(2, len(sd.KeySchema))
	assert.Equal(schema.KeySchema{"Id", "HASH"}, sd.KeySchema[0])
	assert.Equal(1, len(sd.Shards))
	assert.Equal("foo-bar-baz", sd.Shards[0].ShardId)
	assert.Equal("string", sd.Shards[0].ParentShardId)
	assert.Equal("123", sd.Shards[0].SequenceNumberRange.StartingSequenceNumber)
	assert.Equal("456", sd.Shards[0].SequenceNumberRange.EndingSequenceNumber)
	assert.Equal("KEYS_ONLY", sd.StreamViewType)
}

var streamDescription = `{
	"StreamArn": "foo",
	"StreamLabel": "My Stream",
	"TableName": "bar",
	"KeySchema": [
		{"AttributeName": "Id", "KeyType": "HASH"},
		{"AttributeName": "UserId", "KeyType": "RANGE"}
	],
	"LastEvaluatedShardId": "string",
	"Shards": [
		{
			"ShardId": "foo-bar-baz",
			"ParentShardId": "string",
			"SequenceNumberRange": {
				"EndingSequenceNumber": "456",
				"StartingSequenceNumber": "123"
			}
		}
	],
	"StreamViewType": "KEYS_ONLY"
}`
