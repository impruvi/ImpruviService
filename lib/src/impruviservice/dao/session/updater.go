package session

import (
	"../../constants/tablenames"
	"../../util"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
)

func CreateFeedback(sessionNumber int, playerId, drillId string) error {
	session, err := getSession(playerId, sessionNumber)
	if err != nil {
		return err
	}

	found := false
	for _, drill := range session.Drills {
		if drill.DrillId == drillId {
			drill.Feedback = &Feedback{
				VideoUploadDateEpochMillis: util.GetCurrentTimeEpochMillis(),
			}
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Drill %v does not exist for session. playerId %v, sessionNumber: %v", drillId, playerId, sessionNumber))
	}

	return PutSession(session)
}

func CreateSubmission(sessionNumber int, playerId, drillId string) error {
	session, err := getSession(playerId, sessionNumber)
	if err != nil {
		return err
	}

	found := false
	for _, drill := range session.Drills {
		if drill.DrillId == drillId {
			drill.Submission = &Submission{
				VideoUploadDateEpochMillis: util.GetCurrentTimeEpochMillis(),
			}
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Drill %v does not exist for session. playerId %v, sessionNumber: %v", drillId, playerId, sessionNumber))
	}

	return PutSession(session)
}

func CreateSession(session *Session) error {
	latestSessionNumber, err := getLatestSessionNumber(session.PlayerId)
	if err != nil {
		return err
	}
	session.SessionNumber = latestSessionNumber + 1
	return PutSession(session)
}

func PutSession(session *Session) error {
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

func DeleteSession(sessionNumber int, playerId string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			playerIdAttr: {
				S: aws.String(playerId),
			},
			sessionNumberAttr: {
				N: aws.String(strconv.Itoa(sessionNumber)),
			},
		},
		TableName: aws.String(tablenames.DrillsTable),
	}

	_, err := dynamo.DeleteItem(input)
	return err
}
