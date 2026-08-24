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

// MaxAltMessage is the highest valid altMsg value accepted by BuildMessage.
const MaxAltMessage = 4

// BuildMessage formats the gift-exchange notification text for giverName about receiverName.
// altMsg selects a progressively shorter alternate phrasing (0 is the default message;
// 1-4 are shorter still), useful for working around carrier SMS filtering.
func BuildMessage(giverName, receiverName string, altMsg int) string {
	switch altMsg {
	case 1:
		return fmt.Sprintf("Hello %s! Welcome to the family gift exchange! Your gift recipient is %s.", giverName, receiverName)
	case 2:
		return fmt.Sprintf("Hello %s! Your gift recipient is %s.", giverName, receiverName)
	case 3:
		return fmt.Sprintf("Your gift recipient is %s.", receiverName)
	case 4:
		return receiverName
	default:
		return fmt.Sprintf("Hello %s! Welcome to the family gift exchange! Your gift recipient is %s. Merry Christmas!", giverName, receiverName)
	}
}

func send(client messageCreator, from, to, body string) (sid string, status string, err error) {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)

	msg, err := client.CreateMessage(params)
	if err != nil {
		return "", "", err
	}

	if msg.Sid != nil {
		sid = *msg.Sid
	}
	if msg.Status != nil {
		status = *msg.Status
	}

	return sid, status, nil
}

// NewClient returns a Twilio REST client, reading TWILIO_ACCOUNT_SID/TWILIO_AUTH_TOKEN from the environment.
func NewClient() *twilio.RestClient {
	return twilio.NewRestClient()
}

// Send sends a single SMS from `from` to `to` with the given body, returning the
// Twilio message SID and its initial status (typically "queued" or "accepted" —
// final delivery status is only known asynchronously, not from this call).
func Send(client *twilio.RestClient, from, to, body string) (sid string, status string, err error) {
	return send(client.Api, from, to, body)
}
