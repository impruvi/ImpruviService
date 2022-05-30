package dynamoclient

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"sync"
)

var once sync.Once

var instance *dynamodb.DynamoDB

func GetClient() *dynamodb.DynamoDB {

	// create the client just once and share among all calls to dynamoDB
	once.Do(func() {
		var sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		instance = dynamodb.New(sess)
	})

	return instance
}
