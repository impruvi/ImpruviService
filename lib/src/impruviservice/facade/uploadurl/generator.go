package uploadurl

import (
	"fmt"
	"github.com/google/uuid"
	"impruviService/constants/bucketnames"
)

type FileLocation struct {
	BucketName string
	Key        string
	URL        string
}

func GenerateMediaFileLocation(pathPrefix string) *FileLocation {
	randomId := uuid.New().String()
	key := fmt.Sprintf("%s/%s", pathPrefix, randomId)
	bucketName := bucketnames.MediaBucket
	return &FileLocation{
		BucketName: bucketName,
		Key:        key,
		URL:        fmt.Sprintf("https://%s.s3.us-west-2.amazonaws.com/%s", bucketName, key),
	}
}
