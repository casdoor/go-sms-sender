// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package go_sms_sender

import (
	"bytes"
	"encoding/json"
	"github.com/volcengine/volc-sdk-golang/service/sms"
)

type VolcClient struct {
	core       *sms.SMS
	sign       string
	template   string
	smsAccount string
}

func GetVolcClient(accessId, accessKey, sign, region, templateId string, smsAccount []string) *VolcClient {
	client := sms.NewInstance()
	client.Client.SetAccessKey(accessId)
	client.Client.SetSecretKey(accessKey)
	client.SetRegion(region)
	return &VolcClient{
		core:       client,
		sign:       sign,
		template:   templateId,
		smsAccount: smsAccount[0],
	}
}

func (c *VolcClient) SendMessage(param map[string]string, targetPhoneNumber ...string) {
	requestParam, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	if len(targetPhoneNumber) < 1 {
		return
	}
	var phoneNumbers bytes.Buffer
	phoneNumbers.WriteString(targetPhoneNumber[0])
	for _, s := range targetPhoneNumber[1:] {
		phoneNumbers.WriteString(",")
		phoneNumbers.WriteString(s)
	}

	req := &sms.SmsRequest{
		SmsAccount:    c.smsAccount,
		Sign:          c.sign,
		TemplateID:    c.template,
		TemplateParam: string(requestParam),
		PhoneNumbers:  phoneNumbers.String(),
	}
	_, _, err = c.core.Send(req)
	if err != nil {
		panic(err)
	}
}
