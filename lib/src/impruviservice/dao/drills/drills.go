package drills

import (
	"../../awsclients/dynamoclient"
	"../../constants/tablenames"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dynamo = dynamoclient.GetClient()

type Category string
const (
	Dribbling Category = "Dribbling"
	Warmup = "Warmup"
	Shooting = "Shooting"
	Passing = "Passing"
)

type Drill struct {
	DrillId string `json:"drillId"`
	Name string `json:"name"`
	Description string `json:"description"`
	Category Category `json:"category"`
	Equipment []Equipment `json:"equipment"`
	Videos DrillVideos `json:"videos"`
}

type EquipmentType string
const (
	Ball EquipmentType = "Ball"
	Cone = "Cone"
	Space = "Space"
)
type RequirementType string
const (
	Count RequirementType = "Count"
	Dimension = "Dimension"
)
type VideoAngle string
const (
	Front VideoAngle = "Front"
	Side VideoAngle = "Side"
	CloseUp = "CloseUp"
)

type Equipment struct {
	EquipmentType EquipmentType `json:"equipmentType"`
	Requirement interface{} `json:"requirement"`
}

type Requirement struct {
	RequirementType RequirementType `json:"requirementType"`
	Count int `json:"count"`
	Dimensions Dimensions `json:"dimension"`
}

// Dimensions in yards
type Dimensions struct {
	Width int `json:"width"`
	Height int `json:"height"`
}

type DrillVideos struct {
	Front Video `json:"front"`
	Side Video `json:"side"`
	CloseUp Video `json:"closeUp"`
}

type Video struct {
	 FileLocation string `json:"fileLocation"`
}

func GetAllDrills() ([]*Drill, error) {
	var allDrills []*Drill
	scanInput := dynamodb.ScanInput{
		TableName: aws.String(tablenames.DrillsTable),
	}

	result, err := dynamo.Scan(&scanInput)
	if err != nil {
		return nil, fmt.Errorf("error scanning drills. exclusiveStartKey: %v\n" , err)
	}

	drills, err := convertItems(result.Items)
	for _, drill := range drills {
		allDrills = append(allDrills, drill)
	}

	for result.LastEvaluatedKey != nil {
		scanInput = dynamodb.ScanInput{
			ExclusiveStartKey: result.LastEvaluatedKey,
			TableName: aws.String(tablenames.DrillsTable),
		}
		result, err = dynamo.Scan(&scanInput)
		if err != nil {
			return nil, fmt.Errorf("error scanning drills. exclusiveStartKey: %v\n" , err)
		}
		drills, err = convertItems(result.Items)
		for _, drill := range drills {
			allDrills = append(allDrills, drill)
		}
	}

	return allDrills, nil
}

func BatchGetDrills(drillIds []string) (map[string]*Drill, error) {
	if len(drillIds) == 0 {
		return make(map[string]*Drill, 0), nil
	}

	allDrills := make(map[string]*Drill, 0)
	batches := split(drillIds)
	for _, batch := range batches {
		var drillAttrs []map[string]*dynamodb.AttributeValue
		for _, drillId := range batch {
			drillAttrs = append(drillAttrs, map[string]*dynamodb.AttributeValue{
				"drillId": {S: aws.String(drillId)},
			})
		}

		result, err := dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				tablenames.DrillsTable: {Keys: drillAttrs},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error batch getting drills. drillIds: %v: %v", drillIds, err)
		}

		drills, err := convertItems(result.Responses[tablenames.DrillsTable])
		fmt.Printf("Drills: %v\n", drills)
		if err != nil {
			return nil, err
		}
		for _, drill := range drills {
			allDrills[drill.DrillId] = drill
		}
	}
	return allDrills, nil
}


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

func split(drillIds []string) [][]string {
	var batches [][]string
	chunkSize := 100

	for i := 0; i < len(drillIds); i += chunkSize {
		end := i + chunkSize
		if end > len(drillIds) {
			end = len(drillIds)
		}
		batches = append(batches, drillIds[i:end])
	}
	return batches
}
