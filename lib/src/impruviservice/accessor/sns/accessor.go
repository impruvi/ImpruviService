package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"impruviService/clients/sns"
	"log"
)

var systemPhoneNumbers = []string{
	"+14253277259",
	"+18084288591",
}

var snsClient = snsclient.GetClient()

func SendTextToSystem(message string) error {
	for _, phoneNumber := range systemPhoneNumbers {
		err := SendTextMessage(phoneNumber, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendTextMessage(phoneNumber, message string) error {
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
		return err
	} else {
		log.Printf("Publish output: %+v\n", output)
	}
	return nil
}
