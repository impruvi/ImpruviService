package uploadurl

import (
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

type GetVideoThumbnailUploadUrlRequest struct {
	VideoType  VideoType                       `json:"videoType"`
	DemoParams DemoVideoUploadUrlRequestParams `json:"demoParams"`
}

type GetVideoThumbnailUploadUrlResponse struct {
	UploadUrl string `json:"uploadUrl"`
}

func GetVideoThumbnailUploadUrl(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetVideoThumbnailUploadUrlRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	fileLocation, err := getImageFileLocation(request)
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

	return converter.Success(GetVideoThumbnailUploadUrlResponse{
		UploadUrl: uploadUrl,
	})
}

func getImageFileLocation(request GetVideoThumbnailUploadUrlRequest) (*files.FileLocation, error) {
	switch request.VideoType {
	case Demo:
		return files.GetDemoVideoThumbnailFileLocation(
			request.DemoParams.DrillId,
			request.DemoParams.Angle,
		), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported video type: %v\n", request.VideoType))
	}
}
