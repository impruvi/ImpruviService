package uploadurl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"impruviService/clients/s3"
	"impruviService/exceptions"
	"impruviService/facade/uploadurl"
	"log"
	"time"
)

var s3Client = s3client.NewClient()

type GetMediaUploadUrlRequest struct {
	PathPrefix string `json:"pathPrefix"`
}

type GetMediaUploadUrlResponse struct {
	FileLocation string `json:"fileLocation"`
	UploadUrl    string `json:"uploadUrl"`
}

func GetMediaUploadUrl(request *GetMediaUploadUrlRequest) (*GetMediaUploadUrlResponse, error) {
	log.Printf("GetMediaUploadUrlRequest: %+v\n", request)
	err := validateGetMediaUploadUrlRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetMediaUploadUrlRequest: %v\n", err)
		return nil, err
	}

	fileLocation := uploadurl.GenerateMediaFileLocation(request.PathPrefix)

	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(fileLocation.BucketName),
		Key:    aws.String(fileLocation.Key),
	})
	uploadUrl, err := req.Presign(15 * time.Minute)
	if err != nil {
		return nil, err
	}

	return &GetMediaUploadUrlResponse{
		FileLocation: fileLocation.URL,
		UploadUrl:    uploadUrl,
	}, nil
}

func validateGetMediaUploadUrlRequest(request *GetMediaUploadUrlRequest) error {
	if request.PathPrefix == "" {
		return exceptions.InvalidRequestError{Message: "PathPrefix cannot be null/empty"}
	}

	return nil
}
