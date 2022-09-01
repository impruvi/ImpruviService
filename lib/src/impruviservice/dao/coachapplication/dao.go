package coachapplication

import (
	"github.com/google/uuid"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.CoachApplicationsTable,
	reflect.TypeOf(&CoachApplicationDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: applicationIdAttr},
	map[string]dynamo.KeySchema{})

func CreateApplication(application *CoachApplicationDB) error {
	applicationId := uuid.New()
	application.ApplicationId = applicationId.String()
	return mapper.Put(application)
}
