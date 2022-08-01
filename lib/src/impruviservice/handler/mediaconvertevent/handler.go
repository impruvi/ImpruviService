package mediaconvertevent

import (
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	mediaconvertAccessor "impruviService/accessor/mediaconvert"
	drillDao "impruviService/dao/drill"
	drillFacade "impruviService/facade/drill"
	"log"
)

type Detail struct {
	Queue  string `json:"queue"`
	Status string `json:"status"`
	JobId  string `json:"jobId"`
}

type Event struct {
	Detail *Detail `json:"detail"`
}

func HandleMediaConvertEvent(event *Event) (interface{}, error) {
	log.Printf("Event: %+v\n", event)

	outputFileLocation, metadata, err := mediaconvertAccessor.GetJob(event.Detail.JobId, event.Detail.Queue)
	if err != nil {
		log.Printf("Failed to get job metadata for event: %+v. Error: %v\n", event, err)
	}

	log.Printf("OutputFileLocation: %v\n", outputFileLocation)
	log.Printf("Metadata: %v\n", metadata)

	if event.Detail.Status == mediaconvert.JobStatusComplete {
		log.Printf("Media convert job succeeded. Metadata: %+v\n", metadata)
		if metadata.Type == mediaconvertAccessor.FeedbackVideo {

		} else if metadata.Type == mediaconvertAccessor.SubmissionVideo {

		} else if metadata.Type == mediaconvertAccessor.DemoVideo {
			err = handleDemoVideoConversion(outputFileLocation, metadata)
			if err != nil {
				log.Printf("Error while handling demo video conversion: %v\n", err)
			}
			return nil, err
		}
	} else if event.Detail.Status == mediaconvert.JobStatusError {
		log.Printf("Media convert job failed. Metadata: %+v\n", metadata)
		// TODO: send system text message
	}
	return nil, nil
}

func handleDemoVideoConversion(outputFileLocation string, metadata *mediaconvertAccessor.Metadata) error {
	drill, err := drillFacade.GetDrillById(metadata.DemoVideoMedata.DrillId)
	if err != nil {
		log.Printf("Error while getting drill by id: %v. Error %v\n", metadata.DemoVideoMedata.DrillId, err)
	}
	if metadata.DemoVideoMedata.Angle == string(drillDao.FrontAngle) {
		drill.Demos.Front.FileLocation = outputFileLocation
	} else if metadata.DemoVideoMedata.Angle == string(drillDao.SideAngle) {
		drill.Demos.Side.FileLocation = outputFileLocation
	} else {
		drill.Demos.Close.FileLocation = outputFileLocation
	}
	err = drillDao.PutDrill(drill) // can't use drillFacade as that will start media conversion again
	if err != nil {
		log.Printf("Error while updating drill: %v\n", err)
	}
	return err
}
