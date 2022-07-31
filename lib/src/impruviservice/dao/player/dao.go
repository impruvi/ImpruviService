package players

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"impruviService/exceptions"
	"log"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.PlayersTable,
	reflect.TypeOf(&PlayerDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: playerIdAttr},
	map[string]dynamo.KeySchema{
		coachIdIndexName: {PartitionKeyAttributeName: coachIdAttr},
		emailIndexName:   {PartitionKeyAttributeName: emailAttr},
	})

func ListPlayers() ([]*PlayerDB, error) {
	items, err := mapper.Scan()
	if err != nil {
		return nil, err
	}
	return items.([]*PlayerDB), err
}

func GetPlayerById(playerId string) (*PlayerDB, error) {
	item, err := mapper.Get(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(playerId)}})
	if err != nil {
		return nil, err
	}

	return item.(*PlayerDB), nil
}

func GetPlayerByEmail(email string) (*PlayerDB, error) {
	items, err := mapper.Query(
		dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(email)}},
		&dynamo.QueryOptions{
			IndexName: emailIndexName,
		})
	if err != nil {
		return nil, err
	}
	players := items.([]*PlayerDB)
	if len(players) == 0 {
		return nil, exceptions.ResourceNotFoundError{Message: fmt.Sprintf("No player exists with email: %v\n", email)}
	}
	if len(players) > 1 {
		// TODO: we likely want to notify us of this issue. It may be nondeterministic what player this matches
		log.Printf("More than one player exists with the same email")
	}
	return players[0], nil
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
