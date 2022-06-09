package files

import (
	"../constants/bucketnames"
	"fmt"
)

type FileLocation struct {
	BucketName string
	Key        string
	URL        string
}

type Angle string

const (
	Front Angle = "FRONT"
	Side        = "SIDE"
	Close       = "CLOSE"
)

func GetSubmissionVideoFileLocation(playerId string, sessionNumber int, drillId string) *FileLocation {
	key := fmt.Sprintf("%v/%v/%v", playerId, sessionNumber, drillId)
	bucketName := bucketnames.SubmissionsBucket
	return &FileLocation{
		BucketName: bucketName,
		Key:        key,
		URL:        fmt.Sprintf("https://%s.s3.us-west-2.amazonaws.com/%s", bucketName, key),
	}
}

func GetFeedbackVideoFileLocation(playerId string, sessionNumber int, drillId string) *FileLocation {
	key := fmt.Sprintf("%v/%v/%v", playerId, sessionNumber, drillId)
	bucketName := bucketnames.FeedbackBucket
	return &FileLocation{
		BucketName: bucketName,
		Key:        key,
		URL:        fmt.Sprintf("https://%s.s3.us-west-2.amazonaws.com/%s", bucketName, key),
	}
}

func GetDrillVideoFileLocation(drillId string, angle Angle) *FileLocation {
	bucketName := bucketnames.DrillsBucket
	return &FileLocation{
		BucketName: bucketName,
		Key:        drillId,
		URL:        fmt.Sprintf("https://%s.s3.us-west-2.amazonaws.com/%s/%s", bucketName, drillId, angle),
	}
}
