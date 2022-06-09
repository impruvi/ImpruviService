package session

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func convertItems(items []map[string]*dynamodb.AttributeValue) ([]*Session, error) {
	var sessions []*Session
	for _, item := range items {
		var session Session
		err := dynamodbattribute.UnmarshalMap(item, &session)
		if err != nil {
			return nil, fmt.Errorf("error unmashalling session: %v. %v", item, err)
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}
