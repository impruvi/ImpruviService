package notification

import (
	"../awsclients/snsclient"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

var snsClient = snsclient.GetClient()

var PhoneNumbersToNotify = []string{
	"+17202331012",
	"+14253277259",
	"+12067145030",
}

func Notify(message string) {
	for _, phoneNumber := range PhoneNumbersToNotify {
		log.Printf("Notifying: %v\n", phoneNumber)
		output, err := snsClient.Publish(&sns.PublishInput{
			Message:     aws.String(message),
			PhoneNumber: aws.String(phoneNumber),
			MessageAttributes: map[string]*sns.MessageAttributeValue{
				"AWS.MM.SMS.OriginationNumber": {
					DataType:    aws.String("String"),
					StringValue: aws.String("+18444412463"),
				},
				"AWS.SNS.SMS.SMSType": {
					DataType:    aws.String("String"),
					StringValue: aws.String("Transactional"),
				},
			},
		})

		if err != nil {
			log.Printf("Error publishing notification: %v", err)
		} else {
			log.Printf("Publish output: %+v\n", output)
		}
	}
}
