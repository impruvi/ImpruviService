package invitationcodes

import (
<<<<<<< HEAD
=======
	"errors"
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/constants/tablenames"
<<<<<<< HEAD
	"impruviService/exceptions"
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
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
		return nil, exceptions.ResourceNotFoundError{
			Message: fmt.Sprintf("user with invitationCode: %v does not exist", invitationCode),
		}
	}

	return convertItem(result.Item)
}
