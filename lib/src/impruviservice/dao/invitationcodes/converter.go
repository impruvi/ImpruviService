package invitationcodes

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func convertItem(item map[string]*dynamodb.AttributeValue) (*InvitationCodeEntry, error) {
	var invitationCodeEntry InvitationCodeEntry
	err := dynamodbattribute.UnmarshalMap(item, &invitationCodeEntry)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling drill: %v. %v", item, err)
	}
	return &invitationCodeEntry, nil
}
