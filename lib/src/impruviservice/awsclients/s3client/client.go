package s3client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

var once sync.Once

var instance *s3.S3

func NewClient() *s3.S3 {

	// create the client just once and share among all calls to dynamoDB
	once.Do(func() {
		var sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		instance = s3.New(sess, &aws.Config{Region: aws.String("us-west-2")})
	})

	return instance
}
