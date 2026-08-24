package sms

import (
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// messageCreator narrows *twilio.RestClient.Api to just what send needs, for testability.
type messageCreator interface {
	CreateMessage(params *twilioApi.CreateMessageParams) (*twilioApi.ApiV2010Message, error)
}

// BuildMessage formats the gift-exchange notification text for giverName about receiverName.
func BuildMessage(giverName, receiverName string) string {
	return fmt.Sprintf("Hello %s! Your gift recipient is %s. Merry Christmas!", giverName, receiverName)
}

func send(client messageCreator, from, to, body string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)

	_, err := client.CreateMessage(params)
	return err
}

// NewClient returns a Twilio REST client, reading TWILIO_ACCOUNT_SID/TWILIO_AUTH_TOKEN from the environment.
func NewClient() *twilio.RestClient {
	return twilio.NewRestClient()
}

// Send sends a single SMS from `from` to `to` with the given body.
func Send(client *twilio.RestClient, from, to, body string) error {
	return send(client.Api, from, to, body)
}
