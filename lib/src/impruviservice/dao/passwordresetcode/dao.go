package passwordresetcode

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.PasswordResetCodesTable,
	reflect.TypeOf(&PasswordResetCodeEntryDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: emailAttr, RangeKeyAttributeName: creationDateEpochMillisAttr},
	map[string]dynamo.KeySchema{})

func GetResetPasswordCodeEntries(email string) ([]*PasswordResetCodeEntryDB, error) {
	items, err := mapper.Query(dynamo.Key{PartitionKey: &dynamodb.AttributeValue{S: aws.String(email)}}, nil)
	if err != nil {
		return nil, err
	}
	return items.([]*PasswordResetCodeEntryDB), nil
}

func PutResetPasswordCodeEntry(entry *PasswordResetCodeEntryDB) error {
	return mapper.Put(entry)
}
