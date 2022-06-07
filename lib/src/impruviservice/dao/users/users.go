package users

import (
	"../../awsclients/dynamoclient"
	"../../constants/tablenames"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

var dynamo = dynamoclient.GetClient()

type UserType string

const (
	Player UserType = "Player"
	Coach           = "Coach"
)

type User struct {
	UserId string `json:"userId"`
	// CoachId is only present when userType is Player:
	// TODO think about separating players and coaches into separate tables
	CoachUserId    string   `json:"coachUserId"`
	Name           string   `json:"name"`
	UserType       UserType `json:"userType"`
	InvitationCode string   `json:"invitationCode"`
}

func GetPlayersForCoach(coachUserId string) ([]*User, error) {
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.UsersTable),
		IndexName: aws.String("coach-userId-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"coachUserId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(coachUserId),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for _, item := range result.Items {
		var user User
		err = dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			return nil, fmt.Errorf("error unmashalling item: %v. %v", result.Items[0], err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetUserById(userId string) (*User, error) {
	log.Printf("Getting user by userId: %v\n", userId)
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablenames.UsersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {S: aws.String(userId)},
		},
	})

	if err != nil {
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, errors.New(fmt.Sprintf("user with userId: %v does not exist", userId))
	}

	user, err := convertItem(result.Item)
	if err != nil {
		return nil, err
	}
	return user, nil
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

func convertItem(item map[string]*dynamodb.AttributeValue) (*User, error) {
	var user User
	err := dynamodbattribute.UnmarshalMap(item, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling user: %v. %v", item, err)
	}
	return &user, nil

}
