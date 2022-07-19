package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	dynamoClient "impruviService/clients/dynamo"
	"impruviService/exceptions"
	"log"
	"reflect"
)

type Key struct {
	PartitionKey *dynamodb.AttributeValue
	RangeKey     *dynamodb.AttributeValue
}

type KeySchema struct {
	PartitionKeyAttributeName string
	RangeKeyAttributeName     string
}

type QueryOptions struct {
	IndexName string
	Reverse   bool
	Limit     int64
}

type DynamoDBMapper struct {
	tableName              string
	itemType               reflect.Type
	keySchema              KeySchema
	globalSecondaryIndexes map[string]KeySchema
	dynamo                 *dynamodb.DynamoDB
}

func New(
	tableName string,
	itemType reflect.Type,
	keySchema KeySchema,
	globalSecondaryIndexes map[string]KeySchema) *DynamoDBMapper {

	if tableName == "" {
		errorMessage := fmt.Sprintf("table name cannot be empty\n")
		panic(errorMessage)
	}
	if itemType == nil {
		errorMessage := fmt.Sprintf("item type cannot be nil for table: %v\n", tableName)
		panic(errorMessage)
	}
	if keySchema.PartitionKeyAttributeName == "" {
		errorMessage := fmt.Sprintf("partition key cannot be empty for table: %v\n", tableName)
		panic(errorMessage)
	}

	return &DynamoDBMapper{
		tableName:              tableName,
		itemType:               itemType,
		keySchema:              keySchema,
		globalSecondaryIndexes: globalSecondaryIndexes,
		dynamo:                 dynamoClient.GetClient(),
	}
}

func (m *DynamoDBMapper) Query(key Key, options *QueryOptions) (interface{}, error) {
	queryInput := &dynamodb.QueryInput{TableName: aws.String(m.tableName)}

	if options != nil {
		queryInput.ScanIndexForward = aws.Bool(!options.Reverse)

		if options.IndexName != "" {
			if m.globalSecondaryIndexes == nil {
				return nil, exceptions.ResourceNotFoundError{Message: fmt.Sprintf("Index: %v does not exist on table: %v\n", options.IndexName, m.tableName)}
			}
			keySchema, ok := m.globalSecondaryIndexes[options.IndexName]
			if !ok {
				return nil, exceptions.ResourceNotFoundError{
					Message: fmt.Sprintf("Index: %v does not exist on table: %v\n", options.IndexName, m.tableName),
				}
			}
			queryInput.IndexName = aws.String(options.IndexName)
			queryInput.KeyConditions = m.convertToDynamoKeyConditions(key, keySchema)
		} else {
			queryInput.KeyConditions = m.convertToDynamoKeyConditions(key, m.keySchema)
		}

		if options.Limit > 0 {
			queryInput.Limit = aws.Int64(options.Limit)
		}
	} else {
		queryInput.KeyConditions = m.convertToDynamoKeyConditions(key, m.keySchema)
	}

	log.Printf("Query input: %v\n", queryInput)
	result, err := m.dynamo.Query(queryInput)
	if err != nil {
		log.Printf("Error while querying table: %v with key: %v, options: %v. %v\n", m.tableName, key, options, err)
		return nil, err
	}

	return m.convertItems(result.Items)
}

func (m *DynamoDBMapper) Get(key Key) (interface{}, error) {
	result, err := m.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(m.tableName),
		Key:       m.convertToDynamoKey(key),
	})

	if err != nil {
		log.Printf("Error while getting item from table: %v with key: %v. %v\n", m.tableName, key, err)
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, exceptions.ResourceNotFoundError{
			Message: fmt.Sprintf("Item with key: %v does not exist in table: %v\n", key, m.tableName),
		}
	}

	return m.convertItem(result.Item)
}

func (m *DynamoDBMapper) BatchGet(keys []Key) (interface{}, error) {
	itemsConverted := reflect.MakeSlice(reflect.SliceOf(m.itemType), 0, 0)
	if len(keys) == 0 {
		return itemsConverted.Interface(), nil
	}

	batches := m.splitKeys(keys)
	for _, batch := range batches {
		result, err := m.dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				m.tableName: {Keys: m.convertToDynamoKeys(batch)},
			},
		})
		if err != nil {
			log.Printf("Error while batch getting items table: %v with keys: %v. %v\n", m.tableName, keys, err)
			return nil, err
		}

		for _, item := range result.Responses[m.tableName] {
			itemConverted, err := m.convertItem(item)
			if err != nil {
				log.Printf("Error while converting item: %v in batch get for table: %v with keys: %v. %v\n", item, m.tableName, keys, err)
				return nil, err
			}
			itemsConverted = reflect.Append(itemsConverted, reflect.ValueOf(itemConverted))
		}
	}

	return itemsConverted.Interface(), nil
}

func (m *DynamoDBMapper) Put(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(m.tableName),
	}

	_, err = m.dynamo.PutItem(input)
	if err != nil {
		log.Printf("Error while putting item: %v to table: %v. %v", item, m.tableName, err)
	}
	return err
}

func (m *DynamoDBMapper) Delete(key Key) error {
	input := &dynamodb.DeleteItemInput{
		Key:       m.convertToDynamoKey(key),
		TableName: aws.String(m.tableName),
	}

	_, err := m.dynamo.DeleteItem(input)
	if err != nil {
		log.Printf("Error while deleting: %v from table: %v. %v", key, m.tableName, err)
	}
	return err
}

func (m *DynamoDBMapper) Scan() (chan interface{}, chan error, chan bool) {
	itemConvertedChan := make(chan interface{})
	errorChan := make(chan error)
	doneChan := make(chan bool)

	go func() {
		defer func() {
			doneChan <- true
			close(doneChan)
			close(errorChan)
			close(itemConvertedChan)
		}()

		lastEvaluatedKey, err := m.scanPage(itemConvertedChan, nil)
		if err != nil {
			errorChan <- err
			return
		}

		for len(lastEvaluatedKey) > 0 {
			lastEvaluatedKey, err = m.scanPage(itemConvertedChan, lastEvaluatedKey)
			if err != nil {
				errorChan <- err
				return
			}
		}
	}()

	return itemConvertedChan, errorChan, doneChan
}

func (m *DynamoDBMapper) scanPage(itemConvertedChan chan interface{}, lastEvaluatedKey map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	result, err := m.dynamo.Scan(&dynamodb.ScanInput{
		TableName:         aws.String(m.tableName),
		ExclusiveStartKey: lastEvaluatedKey,
	})
	if err != nil {
		log.Printf("Error while scanning table: %v with. %v\n", m.tableName, err)
		return nil, err
	}

	for _, item := range result.Items {
		itemConverted, err := m.convertItem(item)
		if err != nil {
			log.Printf("Error while scanning table: %v with. %v\n", m.tableName, err)
			return nil, err
		}

		itemConvertedChan <- itemConverted
	}

	return result.LastEvaluatedKey, nil
}

// DynamoDB batch get limit is 100 items. Split larger lists into lists of smaller lists
func (m *DynamoDBMapper) splitKeys(keys []Key) [][]Key {
	var batches [][]Key
	chunkSize := 100

	for i := 0; i < len(keys); i += chunkSize {
		end := i + chunkSize
		if end > len(keys) {
			end = len(keys)
		}
		batches = append(batches, keys[i:end])
	}
	return batches
}
