package uploadurl

import (
	"../../awsclients/s3client"
	"../../files"
	"../converter"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

var S3Client = s3client.NewClient()

type VideoType string

const (
	Submission VideoType = "SUBMISSION" // playerId/sessionNumber/drillId
	Feedback             = "FEEDBACK"   // playerId/sessionNumber/drillId
	Demo                 = "DEMO"       // drillId
)

type GetVideoUploadUrlRequest struct {
	VideoType        VideoType                             `json:"videoType"`
	SubmissionParams SubmissionVideoUploadUrlRequestParams `json:"submissionParams"`
	FeedbackParams   FeedbackVideoUploadUrlRequestParams   `json:"feedbackParams"`
	DemoParams       DemoVideoUploadUrlRequestParams       `json:"demoParams"`
}

type SubmissionVideoUploadUrlRequestParams struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
}

type FeedbackVideoUploadUrlRequestParams struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
}

type DemoVideoUploadUrlRequestParams struct {
	DrillId string      `json:"drillId"`
	Angle   files.Angle `json:"angle"` // FRONT/SIDE/CLOSE
}

type GetVideoUploadUrlResponse struct {
	UploadUrl string `json:"uploadUrl"`
}

func GetVideoUploadUrl(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetVideoUploadUrlRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	fileLocation, err := getVideoFileLocation(request)
	if err != nil {
		return converter.BadRequest(err.Error())
	}

	req, _ := S3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(fileLocation.BucketName),
		Key:    aws.String(fileLocation.Key),
	})
	uploadUrl, err := req.Presign(15 * time.Minute)
	if err != nil {
		return converter.InternalServiceError("Error while creating presigned url: %v\n", err)
	}

	return converter.Success(GetVideoUploadUrlResponse{
		UploadUrl: uploadUrl,
	})
}

func getVideoFileLocation(request GetVideoUploadUrlRequest) (*files.FileLocation, error) {
	switch request.VideoType {
	case Feedback:
		return files.GetFeedbackVideoFileLocation(
			request.FeedbackParams.PlayerId,
			request.FeedbackParams.SessionNumber,
			request.FeedbackParams.DrillId,
		), nil
	case Submission:
		return files.GetSubmissionVideoFileLocation(
			request.SubmissionParams.PlayerId,
			request.SubmissionParams.SessionNumber,
			request.SubmissionParams.DrillId,
		), nil
	case Demo:
		return files.GetDemoVideoFileLocation(
			request.DemoParams.DrillId,
			request.DemoParams.Angle,
		), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported video type: %v\n", request.VideoType))
	}
}
