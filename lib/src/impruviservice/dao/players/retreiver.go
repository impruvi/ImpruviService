package players

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/constants/tablenames"
)

func GetPlayerById(playerId string) (*Player, error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablenames.PlayersTable),
		Key: map[string]*dynamodb.AttributeValue{
			playerIdAttr: {S: aws.String(playerId)},
		},
	})

	if err != nil {
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, errors.New(fmt.Sprintf("player with playerId: %v does not exist", playerId))
	}

	return convertItem(result.Item)
}

func GetPlayersForCoach(coachId string) ([]*Player, error) {
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.PlayersTable),
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
