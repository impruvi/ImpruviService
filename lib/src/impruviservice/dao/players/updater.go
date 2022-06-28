package players

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"impruviService/constants/tablenames"
)

func PutPlayer(player *Player) error {
	av, err := dynamodbattribute.MarshalMap(player)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablenames.PlayersTable),
	}

	_, err = dynamo.PutItem(input)
	return err
}
