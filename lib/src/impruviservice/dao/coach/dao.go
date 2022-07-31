package coaches

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"log"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.CoachesTable,
	reflect.TypeOf(&CoachDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: coachIdAttr},
	map[string]dynamo.KeySchema{})

func GetCoachById(coachId string) (*CoachDB, error) {
	item, err := mapper.Get(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(coachId)}})
	if err != nil {
		return nil, err
	}

	return item.(*CoachDB), nil
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
