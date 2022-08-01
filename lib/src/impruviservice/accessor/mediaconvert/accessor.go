package mediaconvert

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	mediaConvertClient "impruviService/clients/mediaconvert"
	"impruviService/constants/mediaconvertqueue"
	"log"
	"strings"
)

type MediaType string

const (
	FeedbackVideo   MediaType = "FEEDBACK_VIDEO"
	SubmissionVideo MediaType = "SUBMISSION_VIDEO"
	DemoVideo       MediaType = "DEMO_VIDEO"
)

type FeedbackVideoMetadata struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
}

type SubmissionVideoMetadata struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
}

type DemoVideoMedata struct {
	DrillId string `json:"drillId"`
	Angle   string `json:"angle"`
}

type Metadata struct {
	Type                    MediaType               `json:"type"`
	FeedbackVideoMetadata   FeedbackVideoMetadata   `json:"feedbackVideoMetadata"`
	SubmissionVideoMetadata SubmissionVideoMetadata `json:"submissionVideoMetadata"`
	DemoVideoMedata         DemoVideoMedata         `json:"demoVideoMedata"`
}

const metadataKey = "metadata"

const maxBitRate = 2000000
const frameRate = 30

var client = mediaConvertClient.NewClient()

func StartJob(inputFileLocation string, metadata *Metadata) error {
	outputS3Location := getOutputS3Location(inputFileLocation)

	log.Printf("Starting job with inputFileLocation: %v. outputFileLocation: %v\n", inputFileLocation, outputS3Location)

	bytes, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Error while serializing start job metadata: %+v. error: %v\n", metadata, err)
	}

	out, err := client.CreateJob(&mediaconvert.CreateJobInput{
		Queue:    aws.String(mediaconvertqueue.MediaconvertQueue),
		Role:     aws.String("arn:aws:iam::522042996447:role/service-role/MediaConvert_Default_Role"),
		Settings: getJobSettings(inputFileLocation, outputS3Location),
		UserMetadata: map[string]*string{
			metadataKey: aws.String(string(bytes)),
		},
	})

	if err != nil {
		log.Printf("Error while starting media convert job with input file location: %v. %v", inputFileLocation, err)
		return err
	}

	bytes, err = json.Marshal(out)
	if err != nil {
		log.Printf("Could not serialize output for logging: %v\n", err)
	} else {
		log.Printf("Start job output: %v\n", string(bytes))
	}

	return nil
}

func GetJob(jobId, queue string) (string, *Metadata, error) {
	job, err := client.GetJob(&mediaconvert.GetJobInput{Id: aws.String(jobId)})
	if err != nil {
		log.Printf("Error while getting job with id: %v. queue: %v. Error: %v\n", jobId, queue, err)
	}
	if *job.Job.Queue != queue {
		log.Printf("Queues do not match. Expected queue: %v. Actual queue: %v\n", queue, job.Job.Queue)
		return "", nil, err
	}

	metadataJSON := job.Job.UserMetadata[metadataKey]
	var metadata Metadata
	err = json.Unmarshal([]byte(*metadataJSON), &metadata)
	if err != nil {
		log.Printf("Failed to deserialize media convert metadataJSON: %v. error: %v\n", metadataJSON, err)
		return "", nil, err
	}

	outputFileLocation := getOutputFileLocation(*job.Job.Settings.Inputs[0].FileInput)

	return outputFileLocation, &metadata, nil
}

func getJobSettings(inputFileLocation, outputS3Location string) *mediaconvert.JobSettings {
	return &mediaconvert.JobSettings{
		Inputs:       getInputs(inputFileLocation),
		OutputGroups: getOutputGroups(outputS3Location),
		TimecodeConfig: &mediaconvert.TimecodeConfig{
			Source: aws.String(mediaconvert.InputTimecodeSourceZerobased),
		},
	}
}

func getOutputS3Location(inputFileLocation string) string {
	outputLocation := strings.Replace(inputFileLocation, "https", "s3", 1)
	outputLocation = strings.Replace(outputLocation, ".s3.us-west-2.amazonaws.com", "", 1)
	return outputLocation + "-transcoded"
}

func getOutputFileLocation(inputFileLocation string) string {
	return inputFileLocation + "-transcoded.mp4"
}

func getInputs(inputFileLocation string) []*mediaconvert.Input {
	return []*mediaconvert.Input{
		{
			AudioSelectors: map[string]*mediaconvert.AudioSelector{
				"Audio Selector 1": {
					DefaultSelection: aws.String(mediaconvert.AudioDefaultSelectionDefault),
				},
			},
			FileInput:      aws.String(inputFileLocation),
			TimecodeSource: aws.String(mediaconvert.InputTimecodeSourceZerobased),
			VideoSelector: &mediaconvert.VideoSelector{
				Rotate: aws.String(mediaconvert.InputRotateAuto),
			},
		},
	}
}

func getOutputGroups(outputS3Location string) []*mediaconvert.OutputGroup {
	return []*mediaconvert.OutputGroup{
		{
			Name: aws.String("File group"),
			OutputGroupSettings: &mediaconvert.OutputGroupSettings{
				FileGroupSettings: &mediaconvert.FileGroupSettings{
					Destination: aws.String(outputS3Location),
				},
				Type: aws.String(mediaconvert.OutputGroupTypeFileGroupSettings),
			},
			Outputs: []*mediaconvert.Output{
				{
					AudioDescriptions: []*mediaconvert.AudioDescription{
						{
							AudioSourceName: aws.String("Audio Selector 1"),
							CodecSettings: &mediaconvert.AudioCodecSettings{
								AacSettings: &mediaconvert.AacSettings{
									Bitrate:    aws.Int64(96000),
									CodingMode: aws.String(mediaconvert.AacCodingModeCodingMode20),
									SampleRate: aws.Int64(44100),
								},
								Codec: aws.String(mediaconvert.AudioCodecAac),
							},
						},
					},
					ContainerSettings: &mediaconvert.ContainerSettings{
						Container: aws.String(mediaconvert.ContainerTypeMp4),
					},
					VideoDescription: &mediaconvert.VideoDescription{
						Height: aws.Int64(1024),
						Width:  aws.Int64(576),
						CodecSettings: &mediaconvert.VideoCodecSettings{
							Codec: aws.String(mediaconvert.VideoCodecH264),
							H264Settings: &mediaconvert.H264Settings{
								FramerateDenominator: aws.Int64(1),
								MaxBitrate:           aws.Int64(maxBitRate),
								FramerateControl:     aws.String(mediaconvert.H264FramerateControlSpecified),
								RateControlMode:      aws.String(mediaconvert.H264RateControlModeQvbr),
								FramerateNumerator:   aws.Int64(frameRate),
								SceneChangeDetect:    aws.String(mediaconvert.H264SceneChangeDetectTransitionDetection),
								QualityTuningLevel:   aws.String(mediaconvert.H264QualityTuningLevelMultiPassHq),
							},
						},
					},
				},
			},
		},
	}
}
