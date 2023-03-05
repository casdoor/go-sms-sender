// Copyright 2023 The Casdoor Authors. All Rights Reserved.
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
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type SubmailClient struct {
	api       string
	appid     string
	signature string
	project   string
}

type SubmailResult struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func buildSubmailPostdata(param map[string]string, appid string, signature string, project string, targetPhoneNumber []string) map[string]string {
	multi := make([]map[string]interface{}, 0, 32)

	for _, phoneNumber := range targetPhoneNumber[0:] {
		multi = append(multi, map[string]interface{}{
			"to":   phoneNumber,
			"vars": param,
		})
	}

	m, _ := json.Marshal(multi)
	postdata := make(map[string]string)
	postdata["appid"] = appid
	postdata["signature"] = signature
	postdata["project"] = project
	postdata["multi"] = string(m)
	return postdata
}

func GetSubmailClient(appid string, signature string, project string) (*SubmailClient, error) {
	submailClient := &SubmailClient{
		api:       "https://api-v4.mysubmail.com/sms/multixsend",
		appid:     appid,
		signature: signature,
		project:   project,
	}
	return submailClient, nil
}

func (c *SubmailClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	postdata := buildSubmailPostdata(param, c.appid, c.signature, c.project, targetPhoneNumber)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range postdata {
		_ = writer.WriteField(key, val)
	}
	contentType := writer.FormDataContentType()
	writer.Close()

	resp, err := http.Post(c.api, contentType, body)
	if err != nil {
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return handleSubmailResult(result)
}

func handleSubmailResult(result []byte) error {
	var submailSuccessResult []SubmailResult
	err := json.Unmarshal(result, &submailSuccessResult)
	if err != nil {
		var submailErrorResult SubmailResult
		err := json.Unmarshal(result, &submailErrorResult)
		if err != nil {
			return err
		}
		return fmt.Errorf(submailErrorResult.Msg)
	}

	errMsg := ""
	for _, submailResult := range submailSuccessResult {
		if submailResult.Status != "success" {
			errMsg = fmt.Sprintf("%s %s", errMsg, submailResult.Msg)
		}
	}
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}

	return err
}
