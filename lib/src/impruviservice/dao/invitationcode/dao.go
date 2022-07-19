package invitationcodes

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.InvitationCodesTable,
	reflect.TypeOf(&InvitationCodeEntryDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: invitationCodeAttr},
	map[string]dynamo.KeySchema{})

func GetInvitationCodeEntry(invitationCode string) (*InvitationCodeEntryDB, error) {
	item, err := mapper.Get(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(invitationCode)}})
	if err != nil {
		return nil, err
	}
	return item.(*InvitationCodeEntryDB), nil
}
