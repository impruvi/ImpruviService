package expo

import (
	"fmt"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func SendPushNotification(title string, body string, expoPushToken string) {
	// To check the token is valid
	pushToken, err := expo.NewExponentPushToken(expoPushToken)
	if err != nil {
		panic(err)
	}

	// Create a new Expo SDK client
	client := expo.NewPushClient(nil)

	// Publish message
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

	// Check errors
	if err != nil {
		panic(err)
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
	}
}
