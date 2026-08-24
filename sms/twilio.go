package sms

import (
	"fmt"
	"regexp"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

var nanpNumber = regexp.MustCompile(`^\+1(\d{3})(\d{3})(\d{4})$`)

// FormatPhoneNumber formats a NANP E.164 number (+1XXXXXXXXXX) as "+1 XXX XXX XXXX"
// for human-readable display. Numbers that don't match this shape are returned unchanged.
func FormatPhoneNumber(number string) string {
	matches := nanpNumber.FindStringSubmatch(number)
	if matches == nil {
		return number
	}

	return fmt.Sprintf("+1 %s %s %s", matches[1], matches[2], matches[3])
}

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
