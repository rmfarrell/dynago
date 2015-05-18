package dynago_test

import (
	"fmt"
	"github.com/crast/dynago"
)

var endpoint, accessKey, secretKey, table string
var client dynago.Client

func Example_query() {
	client := dynago.NewClient(endpoint, accessKey, secretKey)

	query := client.Query(table).
		FilterExpression("NumViews > :views").
		Param(":views", 50).
		Desc()

	result, err := query.Execute()
	if err != nil {
		// do something
	}
	for _, row := range result.Items {
		fmt.Printf("Name: %s, Views: %d", row["Name"], row["NumViews"])
	}
}

func Example_atomicUpdateItem() {
	key := dynago.HashKey("id", 12345)
	result, err := client.UpdateItem("products", key).
		ReturnValues(dynago.ReturnUpdatedNew).
		UpdateExpression("SET SoldCount = SoldCount + :numSold").
		Param(":numSold", 5).
		Execute()
	if err != nil {
		// TODO error handling
	}
	fmt.Printf("We have now sold %d frobbers", result.Attributes["SoldCount"])
}
