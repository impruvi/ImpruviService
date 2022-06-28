package drills

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/constants/tablenames"
)

func GetDrillsForCoach(coachId string) ([]*Drill, error) {
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.DrillsTable),
		IndexName: aws.String(coachIdIndexName),
		KeyConditions: map[string]*dynamodb.Condition{
			coachIdAttr: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(coachId),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return convertItems(result.Items)
}

func BatchGetDrills(drillIds []string) (map[string]*Drill, error) {
	if len(drillIds) == 0 {
		return make(map[string]*Drill, 0), nil
	}

	allDrills := make(map[string]*Drill, 0)
	batches := split(drillIds)
	for _, batch := range batches {
		var drillAttrs []map[string]*dynamodb.AttributeValue
		for _, drillId := range batch {
			drillAttrs = append(drillAttrs, map[string]*dynamodb.AttributeValue{
				drillIdAttr: {S: aws.String(drillId)},
			})
		}

		result, err := dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				tablenames.DrillsTable: {Keys: drillAttrs},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error batch getting drills. drillIds: %v: %v", drillIds, err)
		}

		drills, err := convertItems(result.Responses[tablenames.DrillsTable])
		if err != nil {
			return nil, err
		}
		for _, drill := range drills {
			allDrills[drill.DrillId] = drill
		}
	}
	return allDrills, nil
}

// DynamoDB batch get limit is 100 items. Split larger lists into lists of smaller lists
func split(drillIds []string) [][]string {
	var batches [][]string
	chunkSize := 100

	for i := 0; i < len(drillIds); i += chunkSize {
		end := i + chunkSize
		if end > len(drillIds) {
			end = len(drillIds)
		}
		batches = append(batches, drillIds[i:end])
	}
	return batches
}
