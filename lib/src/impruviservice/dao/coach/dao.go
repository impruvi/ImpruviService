package coaches

import (
	"errors"
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
	tablenames.CoachesTable,
	reflect.TypeOf(&CoachDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: coachIdAttr},
	map[string]dynamo.KeySchema{
		slugIndexName: {PartitionKeyAttributeName: slugAttr},
	})

func GetCoachById(coachId string) (*CoachDB, error) {
	item, err := mapper.Get(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(coachId)}})
	if err != nil {
		return nil, err
	}

	return item.(*CoachDB), nil
}

func GetCoachBySlug(slug string) (*CoachDB, error) {
	items, err := mapper.Query(
		dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(slug)}},
		&dynamo.QueryOptions{IndexName: slugIndexName})

	if err != nil {
		return nil, err
	}
	coaches := items.([]*CoachDB)
	if len(coaches) == 0 {
		return nil, exceptions.ResourceNotFoundError{Message: fmt.Sprintf("No coach with slug: %v", slug)}
	}
	if len(coaches) > 1 {
		return nil, errors.New(fmt.Sprintf("More than 1 coach with slug: %v\n", slug))
	}

	return coaches[0], nil
}

func ListCoaches() ([]*CoachDB, error) {
	items, err := mapper.Scan()
	if err != nil {
		return nil, err
	}
	log.Printf("Items: %+v\n", items)

	return items.([]*CoachDB), nil
}

func PutCoach(coach *CoachDB) error {
	return mapper.Put(coach)
}
