package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"reflect"
)

func (m *DynamoDBMapper) convertItems(items []map[string]*dynamodb.AttributeValue) (interface{}, error) {
	itemsConverted := reflect.MakeSlice(reflect.SliceOf(m.itemType), 0, 0)
	for _, item := range items {
		itemConverted, err := m.convertItem(item)
		if err != nil {
			return nil, err
		}

		itemsConverted = reflect.Append(itemsConverted, reflect.ValueOf(itemConverted))
	}
	return itemsConverted.Interface(), nil
}

func (m *DynamoDBMapper) convertItem(item map[string]*dynamodb.AttributeValue) (interface{}, error) {
	itemConverted := reflect.New(m.itemType)

	err := dynamodbattribute.UnmarshalMap(item, itemConverted.Interface())
	if err != nil {
		return nil, fmt.Errorf("error unmashalling user: %v. %v", item, err)
	}

	return itemConverted.Elem().Interface(), nil
}

func (m *DynamoDBMapper) convertToDynamoKeyConditions(key Key, keySchema KeySchema) map[string]*dynamodb.Condition {
	keyConditions := map[string]*dynamodb.Condition{
		keySchema.PartitionKeyAttributeName: {
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dynamodb.AttributeValue{
				key.PartitionKey,
			},
		},
	}

	if keySchema.RangeKeyAttributeName != "" && key.RangeKey != nil {
		keyConditions[keySchema.RangeKeyAttributeName] = &dynamodb.Condition{
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dynamodb.AttributeValue{
				key.RangeKey,
			},
		}
	}

	return keyConditions
}

func (m *DynamoDBMapper) convertToDynamoKeys(keys []Key) []map[string]*dynamodb.AttributeValue {
	var keyDBs []map[string]*dynamodb.AttributeValue
	for _, key := range keys {
		keyDBs = append(keyDBs, m.convertToDynamoKey(key))
	}
	return keyDBs
}

func (m *DynamoDBMapper) convertToDynamoKey(key Key) map[string]*dynamodb.AttributeValue {
	keyDB := map[string]*dynamodb.AttributeValue{
		m.keySchema.PartitionKeyAttributeName: key.PartitionKey,
	}
	if m.keySchema.RangeKeyAttributeName != "" && key.RangeKey != nil {
		keyDB[m.keySchema.RangeKeyAttributeName] = key.RangeKey
	}

	return keyDB
}
