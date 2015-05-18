package dynago_test

import (
	"fmt"
	"github.com/underarmour/dynago"
	"github.com/underarmour/dynago/schema"
)

func ExampleClient_CreateTable_basic(client *dynago.Client) {
	// NewCreateRequest creates a table with simple defaults.
	// You can use chaining to set the hash and range keys.
	table1 := schema.NewCreateRequest("TableName").
		HashKey("UserId", schema.Number).
		RangeKey("Date", schema.String)

	table1.ProvisionedThroughput.ReadCapacityUnits = 45
	client.CreateTable(table1)
}

func ExampleClient_CreateTable_full(client *dynago.Client) {
	// Most of the time we don't need the full syntax for making create requests
	// It's shown here mostly for purpose of documentation
	req := &schema.CreateRequest{
		TableName: "PersonalPages",
		AttributeDefinitions: []schema.AttributeDefinition{
			{"UserId", schema.Number},
			{"Title", schema.String},
		},
		KeySchema: []schema.KeySchema{
			{"UserId", schema.HashKey},
			{"Title", schema.RangeKey},
		},
		ProvisionedThroughput: schema.ProvisionedThroughput{
			ReadCapacityUnits:  45,
			WriteCapacityUnits: 72,
		},
	}
	if response, err := client.CreateTable(req); err == nil {
		fmt.Printf("table created, status %s", response.TableDescription.TableStatus)
	}
}

func ExampleClient_Query(client *dynago.Client) {
	result, err := client.Query("table").
		FilterExpression("Foo > :val").
		Param(":val", 45).
		Execute()

	if err == nil {
		for _, row := range result.Items {
			fmt.Printf("ID: %s, Foo: %d", row["Id"], row["Foo"])
		}
	}
}

func ExampleClient_PutItem() {

}
