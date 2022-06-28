package invitationcodes

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/constants/tablenames"
)

func GetInvitationCodeEntry(invitationCode string) (*InvitationCodeEntry, error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablenames.InvitationCodesTable),
		Key: map[string]*dynamodb.AttributeValue{
			invitationCodeAttr: {S: aws.String(invitationCode)},
		},
	})

	if err != nil {
		return nil, err
	}

	if result == nil || result.Item == nil {
		return nil, errors.New(fmt.Sprintf("user with invitationCode: %v does not exist", invitationCode))
	}

	return convertItem(result.Item)
}
