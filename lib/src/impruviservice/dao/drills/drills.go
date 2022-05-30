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

type Drill struct {
	DrillId string `json:"drillId"`
	Name string `json:"name"`
	Description string `json:"description"`
	Equipment []Equipment
	Videos []Video
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
type DimensionUnit string
const (
	Feet DimensionUnit = "Feet"
	Yard = "Yard"
)

type Equipment struct {
	EquipmentType EquipmentType `json:"equipmentType"`
	Requirement interface{}
}

type Requirement struct {
	RequirementType RequirementType `json:"requirementType"`
	Count int `json:"count"`
	Dimensions Dimensions `json:"dimension"`
}

type Dimensions struct {
	Width int `json:"width"`
	Height int `json:"height"`
	Unit DimensionUnit `json:"unit"`
}

type Video struct {
	 Angle string `json:"angle"`
	 FileLocation string `json:"fileLocation"`
}

func GetAllDrills(n int) ([]*Drill, error) {
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
		allDrills = append(allDrills, &drill)
		if len(allDrills) >= n {
			return allDrills, nil
		}
	}

	for len(allDrills) < n {
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
			allDrills = append(allDrills, &drill)
			if len(allDrills) >= n {
				return allDrills, nil
			}
		}
	}

	return allDrills, nil
}

func convertItems(items []map[string]*dynamodb.AttributeValue) ([]Drill, error) {
	var drills []Drill
	for _, item := range items {
		var drill Drill
		err := dynamodbattribute.UnmarshalMap(item, &drill)
		if err != nil {
			return nil, fmt.Errorf("error unmashalling drill: %v. %v", item, err)
		}
		drills = append(drills, drill)
	}
	return drills, nil
}
