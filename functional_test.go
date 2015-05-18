package dynago

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/crast/dynago/schema"
)

var state struct {
	client *Client
}

func setUp(t *testing.T) (*assert.Assertions, *Client) {
	if testing.Short() {
		t.SkipNow()
	}
	if state.client == nil {
		endpoint := os.Getenv("DYNAGO_TEST_ENDPOINT")
		if endpoint == "" {
			t.SkipNow()
		}
		state.client = NewClientExecutor(NewAwsExecutor(endpoint, "us-east-1", "AKIAEXAMPLE", "SECRETEXAMPLE"))
		makeTables(t, state.client)
	}
	return assert.New(t), state.client
}

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
	assert, client := setUp(t)
	putResp, err := client.PutItem("Person", Document{"Id": 42, "Name": "Bob"}).Execute()
	assert.NoError(err)
	assert.Nil(putResp)

	response, err := client.GetItem("Person", HashKey("Id", 42)).Execute()
	assert.Equal("Bob", response.Item["Name"])
	assert.IsType(Number("1"), response.Item["Id"])
	assert.Equal(Number("42"), response.Item["Id"])
}
