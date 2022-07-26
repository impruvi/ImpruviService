package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"reflect"
	"strconv"
)

var mapper = dynamo.New(
	tablenames.SessionsTable,
	reflect.TypeOf(&SessionDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: playerIdAttr, RangeKeyAttributeName: sessionNumberAttr},
	map[string]dynamo.KeySchema{})

func GetSessions(playerId string) ([]*SessionDB, error) {
	items, err := mapper.Query(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(playerId)}}, nil)
	if err != nil {
		return nil, err
	}
	return items.([]*SessionDB), nil
}

func GetSession(playerId string, sessionNumber int) (*SessionDB, error) {
	items, err := mapper.Get(dynamo.Key{
		PartitionKey: &dynamodb.AttributeValue{S: aws.String(playerId)},
		RangeKey:     &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(sessionNumber))}})

	if err != nil {
		return nil, err
	}
	return items.(*SessionDB), nil
}

func GetLatestSessionNumber(playerId string) (int, error) {
	items, err := mapper.Query(
		dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(playerId)}},
		&dynamo.QueryOptions{
			Reverse: true,
			Limit:   1,
		})
	if err != nil {
		return -1, err
	}
	sessions := items.([]*SessionDB)
	if len(sessions) == 0 {
		return 0, nil
	}
	return sessions[0].SessionNumber, nil
}

func PutSession(session *SessionDB) error {
	return mapper.Put(session)
}

func DeleteSession(sessionNumber int, playerId string) error {
	return mapper.Delete(dynamo.Key{
		PartitionKey: &dynamodb.AttributeValue{S: aws.String(playerId)},
		RangeKey:     &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(sessionNumber))}})
}
