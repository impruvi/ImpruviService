package coaches

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func convertItem(item map[string]*dynamodb.AttributeValue) (*Coach, error) {
	var user Coach
	err := dynamodbattribute.UnmarshalMap(item, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling user: %v. %v", item, err)
	}
	return &user, nil

}
