package coaches

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/constants/tablenames"
)

func GetCoachById(coachId string) (*Coach, error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablenames.CoachesTable),
		Key: map[string]*dynamodb.AttributeValue{
			coachIdAttr: {S: aws.String(coachId)},
		},
	})

	if err != nil {
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, errors.New(fmt.Sprintf("coach with coachId: %v does not exist", coachId))
	}

	return convertItem(result.Item)
}
