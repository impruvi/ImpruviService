package expo

import (
	"errors"
	"fmt"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	expoClient "impruviService/clients/expo"
	"log"
)

var client = expoClient.NewClient()

func SendPushNotification(title string, body string, expoPushToken string) error {
	log.Printf("Publishing notification with tite: %v, body: %v. expoPushToken: %v\n", title, body, expoPushToken)

	// To check the token is valid
	pushToken, err := expo.NewExponentPushToken(expoPushToken)
	if err != nil {
		log.Printf("Failed to create new exponent push token for: %v. error: %v\n", expoPushToken, err)
		return err
	}

	response, err := client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     body,
			Data:     map[string]string{"withSome": "data"},
			Sound:    "default",
			Title:    title,
			Priority: expo.DefaultPriority,
		},
	)

	if err != nil {
		log.Printf("Failed to publish notification: %v\n", err)
		return err
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		log.Printf(fmt.Sprintf("%v failed", response.PushMessage.To))
		return errors.New(fmt.Sprintf("%v failed", response.PushMessage.To))
	}

	return nil
}
