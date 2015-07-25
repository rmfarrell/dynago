package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixEndpointUrl(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("http://dynamodb.fake.com/", FixEndpointUrl("http://dynamodb.fake.com"))
	assert.Equal("http://dynamodb.fake.com/", FixEndpointUrl("http://dynamodb.fake.com/"))
	assert.Equal("https://dynamodb.fake.com:443/", FixEndpointUrl("https://dynamodb.fake.com"))
	assert.Equal("https://dynamodb.fake.com:443/foo", FixEndpointUrl("https://dynamodb.fake.com/foo"))
}
