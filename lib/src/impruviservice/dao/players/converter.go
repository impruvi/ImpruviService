package players

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


func convertItems(items []map[string]*dynamodb.AttributeValue) ([]*Player, error) {
	var players []*Player
	for _, item := range items {
		player, err := convertItem(item)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

func convertItem(item map[string]*dynamodb.AttributeValue) (*Player, error) {
	var user Player
	err := dynamodbattribute.UnmarshalMap(item, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling user: %v. %v", item, err)
	}
	return &user, nil

}