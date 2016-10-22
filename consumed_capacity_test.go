package dynago_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rmfarrell/dynago"
)

var consumedOk = 0
var indexedBobKey = dynago.Document{"Id": "Hello", "UserId": "bob"}

func consumedCapacitySetup(t *testing.T) (*assert.Assertions, *dynago.Client) {
	assert, client := funcTest.setUp(t)
	if consumedOk == 0 {
		r, err := client.GetItem("Person", dynago.HashKey("Id", 123)).ReturnConsumedCapacity(dynago.CapacityIndexes).Execute()
		if err != nil {
			t.Fatal(err)
		}
		if r.ConsumedCapacity == nil {
			consumedOk = 2
			t.Skip(
				"We believe your server does not support consumed capacity responses.",
				"Consumed Capacity responses are only available from DynamoDB production, not DynamoDB local.",
			)
		} else {
			consumedOk = 1
		}
	}
	if consumedOk == 2 {
		t.SkipNow()
	}
	return assert, client
}

func TestBatchGet_ConsumedCapacity(t *testing.T) {
	assert, client := consumedCapacitySetup(t)
	r4, err := client.BatchGet().
		Get("Posts", dynago.HashRangeKey("UserId", 42, "Dated", 101)).
		ReturnConsumedCapacity(dynago.CapacityIndexes).
		Execute()
	assert.NoError(err)
	assert.Equal(1, len(r4.ConsumedCapacity))
	tc := r4.ConsumedCapacity.GetTable("Posts")
	assert.Equal(0.5, tc.CapacityUnits)
}

func TestBatchWrite_ConsumedCapacity(t *testing.T) {
	assert, client := consumedCapacitySetup(t)
	r3, err := client.BatchWrite().
		ReturnConsumedCapacity(dynago.CapacityIndexes).
		Put("Person", person(123, "Hello")).
		Delete("Indexed", dynago.Document{"Id": "Nothing", "UserId": "bob"}).
		Execute()

	assert.NoError(err)
	assert.Equal(2, len(r3.ConsumedCapacity))
	tc := r3.ConsumedCapacity.GetTable("Person")
	assert.Equal(1.0, tc.CapacityUnits)
	tc = r3.ConsumedCapacity.GetTable("Indexed")
	assert.Equal(1.0, tc.CapacityUnits)
	client.DeleteItem("Person", dynago.HashKey("Id", 123)).Execute()
}

func TestPutItem_ConsumedCapacity(t *testing.T) {
	assert, client := consumedCapacitySetup(t)
	client.DeleteItem("Indexed", indexedBobKey).Execute()
	r1, err := client.PutItem("Indexed", dynago.Document{"Id": "Hello", "UserId": "bob", "Dated": 123}).
		ReturnConsumedCapacity(dynago.CapacityIndexes).
		Execute()
	assert.NoError(err)
	assert.Equal(3.0, r1.ConsumedCapacity.CapacityUnits)
	assert.Equal(1.0, r1.ConsumedCapacity.Table.CapacityUnits)
	assert.Equal(1, len(r1.ConsumedCapacity.GlobalSecondaryIndexes))
	assert.Equal(1.0, r1.ConsumedCapacity.GlobalSecondaryIndexes["index1"].CapacityUnits)
	assert.Equal(1, len(r1.ConsumedCapacity.LocalSecondaryIndexes))
	assert.Equal(1.0, r1.ConsumedCapacity.LocalSecondaryIndexes["index2"].CapacityUnits)
}

func TestQuery_ConsumedCapacity(t *testing.T) {
	assert, client := consumedCapacitySetup(t)
	client.PutItem("Indexed", dynago.Document{"Id": "Hello", "UserId": "bob", "Dated": 123}).Execute()
	// consumed capacity for queries.
	response, err := client.Query("Indexed").
		IndexName("index2").KeyConditionExpression("UserId = :u").Param(":u", "bob").
		ReturnConsumedCapacity(dynago.CapacityIndexes).
		Execute()
	assert.NoError(err)
	assert.NotNil(response.ConsumedCapacity)
	assert.Equal(0.5, response.ConsumedCapacity.CapacityUnits)
	assert.Equal(0, len(response.ConsumedCapacity.GlobalSecondaryIndexes))
	assert.Equal(1, len(response.ConsumedCapacity.LocalSecondaryIndexes))
	client.DeleteItem("Indexed", indexedBobKey).Execute()
}

func TestScan_ConsumedCapacity(t *testing.T) {
	assert, client := consumedCapacitySetup(t)
	r, err := client.Scan("Posts").ReturnConsumedCapacity(dynago.CapacityIndexes).Execute()
	assert.NoError(err)
	if r.ConsumedCapacity.CapacityUnits < 0.5 {
		t.Fatalf("Capacity units too low for a scan of posts %v", *r.ConsumedCapacity)
	}
}
