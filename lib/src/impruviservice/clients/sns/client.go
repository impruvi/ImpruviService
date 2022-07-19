package snsclient

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"sync"
)

var once sync.Once

var instance *sns.SNS

func GetClient() *sns.SNS {

	// create the client just once
	once.Do(func() {
		var sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		instance = sns.New(sess)
	})

	return instance
}
