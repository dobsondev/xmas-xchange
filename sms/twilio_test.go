package sms

import (
	"errors"
	"testing"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type fakeMessageCreator struct {
	gotParams *twilioApi.CreateMessageParams
	sid       string
	status    string
	err       error
}

func (f *fakeMessageCreator) CreateMessage(params *twilioApi.CreateMessageParams) (*twilioApi.ApiV2010Message, error) {
	if f.err != nil {
		return nil, f.err
	}

	f.gotParams = params
	return &twilioApi.ApiV2010Message{Sid: &f.sid, Status: &f.status}, nil
}

func TestBuildMessage(t *testing.T) {
	got := BuildMessage("Alice", "Bob")
	want := "Hello Alice! Your gift recipient is Bob. Merry Christmas!"
	if got != want {
		t.Errorf("BuildMessage() = %q, want %q", got, want)
	}
}

func TestSend_Success(t *testing.T) {
	fake := &fakeMessageCreator{sid: "SM123", status: "queued"}

	sid, status, err := send(fake, "+15550000001", "+15550000002", "hello")
	if err != nil {
		t.Fatalf("send returned error: %v", err)
	}

	if sid != "SM123" {
		t.Errorf("Expected sid %q, got %q", "SM123", sid)
	}
	if status != "queued" {
		t.Errorf("Expected status %q, got %q", "queued", status)
	}
	if fake.gotParams.To == nil || *fake.gotParams.To != "+15550000002" {
		t.Errorf("Expected To %q, got %v", "+15550000002", fake.gotParams.To)
	}
	if fake.gotParams.From == nil || *fake.gotParams.From != "+15550000001" {
		t.Errorf("Expected From %q, got %v", "+15550000001", fake.gotParams.From)
	}
	if fake.gotParams.Body == nil || *fake.gotParams.Body != "hello" {
		t.Errorf("Expected Body %q, got %v", "hello", fake.gotParams.Body)
	}
}

func TestSend_Error(t *testing.T) {
	fake := &fakeMessageCreator{err: errors.New("boom")}

	_, _, err := send(fake, "+15550000001", "+15550000002", "hello")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestFormatPhoneNumber(t *testing.T) {
	testCases := []struct {
		name   string
		number string
		want   string
	}{
		{"standard NANP number", "+15556667777", "+1 555 666 7777"},
		{"non-NANP number unchanged", "+442071234567", "+442071234567"},
		{"too few digits unchanged", "+1555666777", "+1555666777"},
		{"missing plus unchanged", "15556667777", "15556667777"},
		{"empty string unchanged", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatPhoneNumber(tc.number)
			if got != tc.want {
				t.Errorf("FormatPhoneNumber(%q) = %q, want %q", tc.number, got, tc.want)
			}
		})
	}
}
