package go_sms_sender

import (
	"bytes"
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
)

type BceClient struct {
	sign     string
	template string
	core     *sms.Client
}

func GetBceClient(accessId, accessKey, sign, template string, endpoint []string) (*BceClient, error) {
	if len(endpoint) < 1 {
		return nil, fmt.Errorf("missing parameter: endpoint")
	}

	client, err := sms.NewClient(accessId, accessKey, endpoint[0])
	if err != nil {
		return nil, err
	}

	bceClient := &BceClient{
		sign:     sign,
		template: template,
		core:     client,
	}

	return bceClient, nil
}

func (c *BceClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: msg code")
	}

	phoneNumbers := bytes.Buffer{}
	phoneNumbers.WriteString(targetPhoneNumber[0])
	for _, s := range targetPhoneNumber[1:] {
		phoneNumbers.WriteString(",")
		phoneNumbers.WriteString(s)
	}

	contentMap := make(map[string]interface{})
	contentMap["code"] = code

	sendSmsArgs := &api.SendSmsArgs{
		Mobile:      phoneNumbers.String(),
		SignatureId: c.sign,
		Template:    c.template,
		ContentVar:  contentMap,
	}

	_, err := c.core.SendSms(sendSmsArgs)
	if err != nil {
		return err
	}

	return nil
}
