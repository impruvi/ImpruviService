package coaches

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
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

func PutCoach(coach *CoachDB) error {
	return mapper.Put(coach)
}
