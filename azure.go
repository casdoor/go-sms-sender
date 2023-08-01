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
	"net/http"
	"strings"
)

type ACSClient struct {
	apiKey   string
	endpoint string
	template string
}

func GetACSClient(apiKey string, template string, endpoint []string) (*ACSClient, error) {
	if len(endpoint) < 1 {
		return nil, fmt.Errorf("missing parameter: endpoint")
	}

	acsClient := &ACSClient{
		apiKey:   apiKey,
		endpoint: endpoint[0],
		template: template,
	}

	return acsClient, nil
}

func (a *ACSClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) < 1 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	fromPhoneNumber := targetPhoneNumber[0]
	message := a.template

	for key, value := range param {
		message = strings.ReplaceAll(message, "{"+key+"}", value)
	}

	url := fmt.Sprintf("%s/sms?api-version=2021-03-07", a.endpoint)
	payload := map[string]interface{}{
		"from": fromPhoneNumber,
		"body": message,
	}

	client := &http.Client{}
	for i := 1; i < len(targetPhoneNumber); i++ {
		payload["to"] = targetPhoneNumber[i]

		requestBody, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("error creating request body: %w", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Add("Authorization", "Bearer "+a.apiKey)
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}

		resp.Body.Close()
	}

	return nil
}
