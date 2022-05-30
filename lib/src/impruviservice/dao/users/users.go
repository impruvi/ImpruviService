package users

import (
	"../../awsclients/dynamoclient"
	"../../constants/tablenames"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dynamo = dynamoclient.GetClient()

type User struct {
	UserId string `json:"userId"`
	InvitationCode string `json:"invitationCode"`
}

func GetAllUsers() ([]*User, error) {
	var allUsers []*User
	scanInput := dynamodb.ScanInput{
		TableName: aws.String(tablenames.UsersTable),
	}

	result, err := dynamo.Scan(&scanInput)
	if err != nil {
		return nil, fmt.Errorf("error scanning users. exclusiveStartKey: %v\n" , err)
	}

	users, err := convertItems(result.Items)
	for _, user := range users {
		allUsers = append(allUsers, &user)
	}

	for result.LastEvaluatedKey != nil {
		scanInput = dynamodb.ScanInput{
			ExclusiveStartKey: result.LastEvaluatedKey,
			TableName: aws.String(tablenames.UsersTable),
		}
		result, err = dynamo.Scan(&scanInput)
		if err != nil {
			return nil, fmt.Errorf("error scanning users. exclusiveStartKey: %v\n" , err)
		}
		users, err = convertItems(result.Items)
		for _, user := range users {
			allUsers = append(allUsers, &user)
		}
	}

	return allUsers, nil
}

func GetUserByInvitationCode(invitationCode string) (*User, error) {
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.UsersTable),
		IndexName: aws.String("invitation-code-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"invitationCode": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(invitationCode),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, errors.New(fmt.Sprintf("User with invitationCode %s does not exist", invitationCode))
	}

	var user User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling item: %v. %v", result.Items[0], err)
	}
	return &user, nil
}

func convertItems(items []map[string]*dynamodb.AttributeValue) ([]User, error) {
	var users []User
	for _, item := range items {
		var user User
		err := dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			return nil, fmt.Errorf("error unmashalling user: %v. %v", item, err)
		}
		users = append(users, user)
	}
	return users, nil
}
