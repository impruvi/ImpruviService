package session

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"impruviService/constants/tablenames"
	"log"
	"strconv"
)

func GetSessions(playerId string) ([]*Session, error) {
	log.Printf("Getting sessions for player with playerId: %v\n", playerId)
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.SessionsTable),
		KeyConditions: map[string]*dynamodb.Condition{
			playerIdAttr: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(playerId),
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

func GetSession(playerId string, sessionNumber int) (*Session, error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablenames.SessionsTable),
		Key: map[string]*dynamodb.AttributeValue{
			playerIdAttr:      {S: aws.String(playerId)},
			sessionNumberAttr: {N: aws.String(strconv.Itoa(sessionNumber))},
		},
	})

	if err != nil {
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, errors.New(fmt.Sprintf("session with playerId: %v sessionNumber: %v does not exist\n", playerId, sessionNumber))
	}

	var session Session
	err = dynamodbattribute.UnmarshalMap(result.Item, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func getLatestSessionNumber(playerId string) (int, error) {
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.PlayersTable),
		KeyConditions: map[string]*dynamodb.Condition{
			playerIdAttr: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(playerId),
					},
				},
			},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int64(1),
	})
	if err != nil {
		return -1, err
	}
	if len(result.Items) == 0 {
		return 0, nil
	}

	sessions, err := convertItems(result.Items)
	log.Printf("Sessions for latest session number: %v\n", sessions)
	return sessions[0].SessionNumber, nil
}
