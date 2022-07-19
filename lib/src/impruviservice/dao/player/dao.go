package players

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.PlayersTable,
	reflect.TypeOf(&PlayerDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: playerIdAttr},
	map[string]dynamo.KeySchema{
		coachIdIndexName: {PartitionKeyAttributeName: coachIdAttr},
	})

func GetPlayerById(playerId string) (*PlayerDB, error) {
	item, err := mapper.Get(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(playerId)}})
	if err != nil {
		return nil, err
	}

	return item.(*PlayerDB), nil
}

func GetPlayersForCoach(coachId string) ([]*PlayerDB, error) {
	items, err := mapper.Query(
		dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(coachId)}},
		&dynamo.QueryOptions{
			IndexName: coachIdIndexName,
		})
	if err != nil {
		return nil, err
	}
	return items.([]*PlayerDB), nil
}

func PutPlayer(player *PlayerDB) error {
	return mapper.Put(player)
}
