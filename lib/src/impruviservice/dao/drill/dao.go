package drills

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"log"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.DrillsTable,
	reflect.TypeOf(&DrillDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: drillIdAttr},
	map[string]dynamo.KeySchema{
		coachIdIndexName: {PartitionKeyAttributeName: coachIdAttr},
	})

func GetDrillById(drillId string) (*DrillDB, error) {
	item, err := mapper.Get(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(drillId)}})
	if err != nil {
		return nil, err
	}
	return item.(*DrillDB), nil
}

func GetDrillsForCoach(coachId string) ([]*DrillDB, error) {
	items, err := mapper.Query(
		dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(coachId)}},
		&dynamo.QueryOptions{IndexName: coachIdIndexName})

	if err != nil {
		return nil, err
	}
	return items.([]*DrillDB), nil
}

func BatchGetDrills(drillIds []string) (map[string]*DrillDB, error) {
	keys := make([]dynamo.Key, 0)
	for _, drillId := range drillIds {
		keys = append(keys, dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(drillId)}})
	}
	items, err := mapper.BatchGet(keys)
	if err != nil {
		return nil, err
	}

	drills := items.([]*DrillDB)
	drillsMap := make(map[string]*DrillDB, 0)
	for _, drill := range drills {
		drillsMap[drill.DrillId] = drill
	}
	return drillsMap, nil
}

func CreateDrill(drill *DrillDB) (*DrillDB, error) {
	drillId := uuid.New()
	drill.DrillId = drillId.String()
	err := PutDrill(drill)
	if err != nil {
		return nil, err
	}
	return drill, nil
}

func PutDrill(drill *DrillDB) error {
	return mapper.Put(drill)
}

func DeleteDrill(drillId string) error {
	drill, err := GetDrillById(drillId)
	log.Printf("Deleting drill: %v\n", drill)
	if err != nil {
		log.Printf("Unexpected error while getting drill to delete: %v\n", drill)
		return err
	}
	drill.IsDeleted = true
	return PutDrill(drill)
}
