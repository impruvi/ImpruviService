package coaches

import (
	"../../constants/tablenames"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutCoach(coach *Coach) error {
	av, err := dynamodbattribute.MarshalMap(coach)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablenames.CoachesTable),
	}

	_, err = dynamo.PutItem(input)
	return err
}
