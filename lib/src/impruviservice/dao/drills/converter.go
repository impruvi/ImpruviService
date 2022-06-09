package drills

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func convertItems(items []map[string]*dynamodb.AttributeValue) ([]*Drill, error) {
	var drills []*Drill
	for _, item := range items {
		var drill Drill
		err := dynamodbattribute.UnmarshalMap(item, &drill)
		if err != nil {
			return nil, fmt.Errorf("error unmashalling drill: %v. %v", item, err)
		}
		drills = append(drills, &drill)
	}
	return drills, nil
}
