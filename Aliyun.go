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

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type AliyunClient struct {
	template string
	sign string
	core *dysmsapi.Client
}

func GetAliyunClient(accessId, accessKey, sign, region, template string) *AliyunClient {
	client, err := dysmsapi.NewClientWithAccessKey(region, accessId, accessKey)
	if err != nil {
		panic(err)
	}
	return &AliyunClient{
		template: template,
		core: client,
		sign: sign,
	}
}

func (c *AliyunClient) SendMessage(param map[string]string, targetPhoneNumber ...string) {
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

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumbers.String()
	request.TemplateCode = c.template
	request.TemplateParam = string(requestParam)
	request.SignName = c.sign

	_, err = c.core.SendSms(request)
	if err != nil {
		panic(err)
	}
}
