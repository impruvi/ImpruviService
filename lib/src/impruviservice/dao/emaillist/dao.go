package emaillist

import (
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"impruviService/util"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.EmailListSubscriptionsTable,
	reflect.TypeOf(&EmailListSubscriptionDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: emailAttr},
	map[string]dynamo.KeySchema{})

func CreateSubscription(email string) error {
	currentTime := util.GetCurrentTimeEpochMillis()
	return mapper.Put(&EmailListSubscriptionDB{
		Email:                      email,
		CreationDateEpochMillis:    currentTime,
		LastUpdatedDateEpochMillis: currentTime,
	})
}
