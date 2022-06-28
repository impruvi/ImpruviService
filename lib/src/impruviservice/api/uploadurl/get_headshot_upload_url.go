package uploadurl

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"impruviService/api/converter"
	"impruviService/files"
	"time"
)

type GetHeadshotUploadUrlRequest struct {
	UserType string `json:"userType"`
	UserId   string `json:"userId"`
}

type GetHeadshotUploadUrlResponse struct {
	UploadUrl string `json:"uploadUrl"`
}

func GetHeadshotUploadUrl(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetHeadshotUploadUrlRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	fileLocation, err := files.GetHeadshotFileLocation(
		request.UserType,
		request.UserId,
	), nil
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
