package uploadurl

import (
	"../../awsclients/s3client"
	"../../constants/bucketnames"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"time"
)

var S3Client = s3client.NewClient()

type VideoType string
const (
	Submission VideoType = "Submission"
	Feedback = "Feedback"
)

type GetVideoUploadUrlRequest struct {
	VideoType VideoType `json:"videoType"`
	UserId string `json:"userId"`
	SessionNumber int `json:"sessionNumber"`
	DrillId string `json:"drillId"`
}

type GetVideoUploadUrlResponse struct {
	Url string `json:"url"`
}

func GetVideoUploadUrl(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetVideoUploadUrlRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	req, _ := S3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketnames.SubmissionsBucket),
		Key:    aws.String(fmt.Sprintf("%v/%v/%v", request.UserId, request.SessionNumber, request.DrillId)),
	})
	uploadUrl, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Error while creating presigned url: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	rspBody, err := json.Marshal(GetVideoUploadUrlResponse{
		Url: uploadUrl,
	})
	if err != nil {
		log.Printf("Error while marshalling response: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		Body:       string(rspBody),
		StatusCode: http.StatusAccepted,
	}
}
