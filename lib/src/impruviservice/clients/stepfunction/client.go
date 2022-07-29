package stepfunction

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"sync"
)

var once sync.Once

var instance *sfn.SFN

func GetClient() *sfn.SFN {

	// create the client just once
	once.Do(func() {
		var sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		instance = sfn.New(sess)
	})

	return instance
}
