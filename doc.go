/*
Dynago is a DynamoDB client API for Go.

This attempts to be a really simple, principle of least-surprise client for the DynamoDB API.

Key design tenets of Dynago:

	* Most actions are done via chaining to build filters and conditions
	* all objects are completely safe for passing between goroutines (even queries and the like)
	  ^ this is because there's no shared state
	* To make understanding easier via docs, we use amazon's naming wherever possible.

EXAMPLES

Run a query:

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


Future:
	* Automatic cursor pagination
	* (maybe) Ability to map results into a struct

*/
package dynago
