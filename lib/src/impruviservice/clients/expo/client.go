package expo

import (
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"sync"
)

var once sync.Once

var instance *expo.PushClient

func NewClient() *expo.PushClient {

	// create the client just once and share among all calls
	once.Do(func() {
		instance = expo.NewPushClient(nil)
	})

	return instance
}
