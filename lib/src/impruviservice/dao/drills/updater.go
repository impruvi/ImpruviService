package drills

import (
	"../../constants/tablenames"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

func CreateDrill(drill *Drill) error {
	drillId := uuid.New()
	drill.DrillId = drillId.String()
	return PutDrill(drill)
}

func PutDrill(drill *Drill) error {
	av, err := dynamodbattribute.MarshalMap(drill)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablenames.DrillsTable),
	}

	_, err = dynamo.PutItem(input)
	return err
}

func DeleteDrill(drillId string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			drillIdAttr: {
				N: aws.String(drillId),
			},
		},
		TableName: aws.String(tablenames.DrillsTable),
	}

	_, err := dynamo.DeleteItem(input)
	return err
}
