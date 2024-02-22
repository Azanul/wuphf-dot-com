package twiliosms

import (
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSMS struct {
	client *twilio.RestClient
}

func New() *TwilioSMS {
	return &TwilioSMS{
		client: twilio.NewRestClient(),
	}
}

func (ts *TwilioSMS) Name() string {
	return "twilio sms"
}

func (ts *TwilioSMS) Notify(receiver, body string) (string, error) {
	sender, err := ts.client.NumbersV2.ListHostedNumberOrder(nil)
	if err != nil {
		return "", err
	}
	params := &api.CreateMessageParams{}
	params.SetBody(body)
	params.SetFrom(*sender[0].PhoneNumber)
	params.SetTo(receiver)

	resp, err := ts.client.Api.CreateMessage(params)
	if err != nil {
		return "", err
	} else if resp.Sid == nil {
		return "", nil
	}
	return *resp.Sid, err
}
