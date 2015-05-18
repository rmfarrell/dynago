package dynago

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/underarmour/dynago/schema"
)

type functional struct {
	client *Client
}

func (f *functional) setUp(t *testing.T) (*assert.Assertions, *Client) {
	if testing.Short() {
		t.SkipNow()
	}
	if f.client == nil {
		endpoint := os.Getenv("DYNAGO_TEST_ENDPOINT")
		if endpoint == "" {
			t.SkipNow()
		}
		f.client = NewClientExecutor(NewAwsExecutor(endpoint, "us-east-1", "AKIAEXAMPLE", "SECRETEXAMPLE"))
		makeTables(t, f.client)
	}
	return assert.New(t), f.client
}

var funcTest functional

func makeTables(t *testing.T, client *Client) {
	hashTable := schema.NewCreateRequest("Person").HashKey("Id", schema.Number)
	hashRange := schema.NewCreateRequest("Posts").
		HashKey("UserId", schema.Number).
		RangeKey("Dated", schema.String)

	tables := []*schema.CreateRequest{hashTable, hashRange}
	for _, table := range tables {
		_, err := client.CreateTable(table)
		if err != nil {
			panic(err)
		}
	}
}

func TestGet(t *testing.T) {
	assert, client := funcTest.setUp(t)
	putResp, err := client.PutItem("Person", Document{"Id": 42, "Name": "Bob"}).Execute()
	assert.NoError(err)
	assert.Nil(putResp)

	response, err := client.GetItem("Person", HashKey("Id", 42)).Execute()
	assert.Equal("Bob", response.Item["Name"])
	assert.IsType(Number("1"), response.Item["Id"])
	assert.Equal(Number("42"), response.Item["Id"])
}
