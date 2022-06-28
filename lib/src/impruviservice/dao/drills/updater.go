package drills

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"impruviService/constants/tablenames"
)

func CreateDrill(drill *Drill) (*Drill, error) {
	drillId := uuid.New()
	drill.DrillId = drillId.String()
	err := PutDrill(drill)
	if err != nil {
		return nil, err
	}
	return drill, err
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
				S: aws.String(drillId),
			},
		},
		TableName: aws.String(tablenames.DrillsTable),
	}

	_, err := dynamo.DeleteItem(input)
	return err
}
