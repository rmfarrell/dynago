package dynago_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/underarmour/dynago"
	"github.com/underarmour/dynago/schema"
)

type functional struct {
	client *dynago.Client
}

func (f *functional) setUp(t *testing.T) (*assert.Assertions, *dynago.Client) {
	if testing.Short() {
		t.SkipNow()
	}
	if f.client == nil {
		endpoint := os.Getenv("DYNAGO_TEST_ENDPOINT")
		if endpoint == "" {
			t.SkipNow()
		}
		executor := dynago.NewAwsExecutor(endpoint, "us-east-1", "AKIAEXAMPLE", "SECRETEXAMPLE")
		f.client = dynago.NewClient(executor)
		makeTables(t, f.client)
	}
	return assert.New(t), f.client
}

var funcTest functional

func makeTables(t *testing.T, client *dynago.Client) {
	hashTable := schema.NewCreateRequest("Person").HashKey("Id", schema.Number)
	hashRange := schema.NewCreateRequest("Posts").
		HashKey("UserId", schema.Number).
		RangeKey("Dated", schema.String)

	tables := []*schema.CreateRequest{hashTable, hashRange}
	for _, table := range tables {
		_, err := client.CreateTable(table)
		if err != nil {
			if e, ok := err.(*dynago.Error); ok && e.Type == dynago.ErrorResourceInUse {
				continue
			}
			panic(err)
		}
	}
}

func TestGet(t *testing.T) {
	assert, client := funcTest.setUp(t)
	putResp, err := client.PutItem("Person", person(42, "Bob")).Execute()
	assert.NoError(err)
	assert.Nil(putResp)

	response, err := client.GetItem("Person", dynago.HashKey("Id", 42)).Execute()
	assert.Equal("Bob", response.Item["Name"])
	assert.IsType(dynago.Number("1"), response.Item["Id"])
	assert.Equal(dynago.Number("42"), response.Item["Id"])
}

func TestConditionalPut(t *testing.T) {
	assert, client := funcTest.setUp(t)
	doc := person(45, "Joe")
	doc["Count"] = 94
	client.PutItem("Person", doc).Execute()

	doc["Count"] = 45

	basePut := client.PutItem("Person", doc).
		ConditionExpression("#c > :val").
		Param("#c", "Count")

	_, err := basePut.Param(":val", 100).Execute()

	e := err.(*dynago.Error)
	assert.Equal(dynago.ErrorConditionFailed, e.Type)

	_, err = basePut.Param(":val", 50).Execute()
	assert.NoError(err)
}

func TestBatchWrite(t *testing.T) {
	assert, client := funcTest.setUp(t)
	_, err := client.PutItem("Person", person(4, "ToDelete")).Execute()
	assert.NoError(err)

	p1 := person(1, "Joe")
	p2 := person(2, "Mary")
	p3 := person(3, "Amy")
	_, err = client.BatchWrite().
		Put("Person", p1, p2, p3).
		Delete("Person", dynago.HashKey("Id", 4)).
		Execute()

	assert.NoError(err)

	response, err := client.GetItem("Person", dynago.HashKey("Id", 2)).Execute()
	assert.Equal("Mary", response.Item["Name"])
}

func TestUpdateItemConditional(t *testing.T) {
	assert, client := funcTest.setUp(t)
	_, err := client.PutItem("Person", person(5, "ToUpdate")).Execute()
	assert.NoError(err)
	result, err := client.UpdateItem("Person", dynago.HashKey("Id", 5)).
		UpdateExpression("SET #n = :name").
		Param("#n", "Name").Param(":name", "Bob").
		ReturnValues(dynago.ReturnUpdatedNew).Execute()
	assert.NoError(err)
	assert.NotNil(result)
	assert.NotNil(result.Attributes)
	assert.Equal("Bob", result.Attributes["Name"])
}

func TestTableActions(t *testing.T) {
	tables := []string{"abc", "def", "ghi"}
	assert, client := funcTest.setUp(t)
	for _, name := range tables {
		_, err := client.CreateTable(schema.NewCreateRequest(name).HashKey("Id", schema.Number))
		if e, ok := err.(*dynago.Error); !ok || e.Type != dynago.ErrorResourceInUse {
			assert.NoError(err)
		}
	}
	list, err := client.ListTables().Limit(10).Execute()
	assert.NoError(err)
	assert.NotNil(list)
	assert.True(len(list.TableNames) > len(tables))
	assert.Nil(list.Next())

	// Pagination of tables should work
	list1, err := client.ListTables().Limit(2).Execute()
	assert.NoError(err)
	assert.NotNil(list1.Next())
	assert.Equal(2, len(list1.TableNames))
	assert.Equal(list.TableNames[:2], list1.TableNames[:2])
	list2, err := list1.Next().Execute()
	assert.Equal(list.TableNames[2:4], list2.TableNames[:2])

	for _, name := range tables {
		_, err := client.DeleteTable(name)
		assert.NoError(err)
	}

}

func TestDescribeTable(t *testing.T) {
	assert, client := funcTest.setUp(t)
	response, err := client.DescribeTable("Bogus")
	assert.Error(err)

	response, err = client.DescribeTable("Posts")
	assert.NoError(err)
	assert.Equal("Posts", response.Table.TableName)
	assert.Equal(2, len(response.Table.AttributeDefinitions))
	assert.Equal(0, len(response.Table.GlobalSecondaryIndexes))
}

func person(id int, name string) dynago.Document {
	return dynago.Document{"Id": id, "Name": name, "IncVal": 1}
}
