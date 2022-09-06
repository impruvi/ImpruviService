package coachapplication

import (
	"github.com/google/uuid"
	"impruviService/accessor/dynamo"
	"impruviService/constants/tablenames"
	"impruviService/util"
	"reflect"
)

var mapper = dynamo.New(
	tablenames.CoachApplicationsTable,
	reflect.TypeOf(&CoachApplicationDB{}),
	dynamo.KeySchema{PartitionKeyAttributeName: applicationIdAttr},
	map[string]dynamo.KeySchema{})

func CreateApplication(application *CoachApplicationDB) error {
	currentTime := util.GetCurrentTimeEpochMillis()
	applicationId := uuid.New()
	application.ApplicationId = applicationId.String()
	application.CreationDateEpochMillis = currentTime
	application.LastUpdatedDateEpochMillis = currentTime
	return mapper.Put(application)
}
