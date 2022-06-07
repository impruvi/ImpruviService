package warmup

import (
	"../../awsclients/dynamoclient"
	"../../constants/tablenames"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// HandleWarmupEvent primes the DB connections
func HandleWarmupEvent() {
	fmt.Printf("Handling warm up event")
	dynamo := dynamoclient.GetClient()

	tableNames := []string{
		tablenames.UsersTable,
		tablenames.DrillsTable,
		tablenames.SessionsTable,
	}
	for _, tableName := range tableNames {
		describeTableInput := dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}
		_, err := dynamo.DescribeTable(&describeTableInput)

		if err != nil {
			fmt.Printf("Error while describing table: %s, e: %s", tableName, err)
		}
	}
}
