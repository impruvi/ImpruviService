package mediaconvert

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"sync"
)

var once sync.Once

var instance *mediaconvert.MediaConvert

func NewClient() *mediaconvert.MediaConvert {

	// create the client just once
	once.Do(func() {
		var sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		instance = mediaconvert.New(sess, &aws.Config{
			Region:   aws.String("us-west-2"),
			Endpoint: aws.String("https://mlboolfjb.mediaconvert.us-west-2.amazonaws.com"),
		})
	})

	return instance
}
