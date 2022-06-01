package session

import (
	"../../awsclients/dynamoclient"
	"../../constants/tablenames"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
	"time"
)

var dynamo = dynamoclient.GetClient()

type Session struct {
	UserId string `json:"userId"`
	SessionNumber int `json:"sessionNumber"`
	Drills []*Drill `json:"drills"`
}

type Drill struct {
	DrillId string `json:"drillId"`
	Submission *Submission `json:"submission"`
	Feedback *Feedback `json:"feedback"`
	Tips []string `json:"tips"`
	Repetitions int `json:"repetitions"`
	DurationMinutes int `json:"durationMinutes"`
}

type Submission struct {
	CreationDateEpochMillis int64 `json:"creationDateEpochMillis"`
	FileLocation string `json:"fileLocation"`
}

type Feedback struct {
	CreationDateEpochMillis int64 `json:"creationDateEpochMillis"`
	FileLocation string `json:"fileLocation"`
}


func GetSessions(userId string) ([]*Session, error) {
	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(tablenames.SessionsTable),
		KeyConditions: map[string]*dynamodb.Condition{
			"userId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userId),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return convertItems(result.Items)
}

func CreateFeedback(sessionNumber int, userId, drillId, fileLocation string) error {
	session, err := getSession(userId, sessionNumber)
	if err != nil {
		return err
	}

	found := false
	for _, drill := range session.Drills {
		if drill.DrillId == drillId {
			drill.Feedback = &Feedback{
				CreationDateEpochMillis: time.Now().UnixNano() / int64(time.Millisecond),
				FileLocation:            fileLocation,
			}
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Drill %v does not exist for session. userId %v, sessionNumber: %v", drillId, userId, sessionNumber))
	}

	return putSession(session)
}

func CreateSubmission(sessionNumber int, userId, drillId, fileLocation string) error {
	session, err := getSession(userId, sessionNumber)
	if err != nil {
		return err
	}

	found := false
	for _, drill := range session.Drills {
		if drill.DrillId == drillId {
			drill.Submission = &Submission{
				CreationDateEpochMillis: time.Now().UnixNano() / int64(time.Millisecond),
				FileLocation:            fileLocation,
			}
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Drill %v does not exist for session. userId %v, sessionNumber: %v", drillId, userId, sessionNumber))
	}

	return putSession(session)
}

func UpdateSession(sessionNumber int, userId string, session *Session) error {
	session, err := getSession(userId, sessionNumber)
	if err != nil {
		return err
	}

	return putSession(session)
}

func putSession(session *Session) error {
	av, err := dynamodbattribute.MarshalMap(session)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablenames.SessionsTable),
	}

	_, err = dynamo.PutItem(input)
	return err
}


func getSession(userId string, sessionNumber int) (*Session, error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablenames.SessionsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {S: aws.String(userId)},
			"sessionNumber": {N: aws.String(strconv.Itoa(sessionNumber))},
		},
	})

	if err != nil {
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, errors.New(fmt.Sprintf("session with userId: %v sessionNumber: %v does not exist\n", userId, sessionNumber))
	}

	var session Session
	err = dynamodbattribute.UnmarshalMap(result.Item, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func convertItems(items []map[string]*dynamodb.AttributeValue) ([]*Session, error) {
	var sessions []*Session
	for _, item := range items {
		var session Session
		err := dynamodbattribute.UnmarshalMap(item, &session)
		if err != nil {
			return nil, fmt.Errorf("error unmashalling session: %v. %v", item, err)
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}
