package ses

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"sync"
)

var once sync.Once

var instance *ses.SES

func GetClient() *ses.SES {

	// create the client just once
	once.Do(func() {
		var sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		instance = ses.New(sess)
	})

	return instance
}
