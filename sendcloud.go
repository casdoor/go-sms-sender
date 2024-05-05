// Copyright 2024 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package go_sms_sender

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SendCloudResponseData ResponseData structure holds the response from SendCloud API.
type SendCloudResponseData struct {
	Result     bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Info       string `json:"info"`
}

// SendCloudConfig Config holds the configuration for the SMS sending service.
type SendCloudConfig struct {
	CharSet      string
	Server       string
	SendSMSAPI   string
	MaxReceivers int
}

type SendCloudClient struct {
	SmsUser    string
	SmsKey     string
	TemplateId int
	MsgType    int
	Config     SendCloudConfig
}

func GetSendCloudClient(smsUser string, smsKey string, template string) (*SendCloudClient, error) {
	templateId, err := strconv.Atoi(template)
	if err != nil {
		return nil, fmt.Errorf("template id should be number")
	}

	msgType, err := strconv.Atoi(template)
	if err != nil {
		return nil, fmt.Errorf("msgType id should be number")
	}

	return &SendCloudClient{
		SmsUser:    smsUser,
		SmsKey:     smsKey,
		TemplateId: templateId,
		MsgType:    msgType,
		Config: SendCloudConfig{
			CharSet:      "utf-8",
			Server:       "https://api.sendcloud.net",
			SendSMSAPI:   "/smsapi/send",
			MaxReceivers: 100,
		},
	}, nil
}

// SendMessage sends an SMS using SendCloud API.
func (client *SendCloudClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if err := client.validateSendCloudSms(targetPhoneNumber); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	if err := validateConfig(client.Config); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	params, err := client.prepareParams(param, targetPhoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	signature := calculateSignature(params, client.SmsKey)
	params.Set("signature", signature)

	resp, err := http.PostForm(client.Config.SendSMSAPI, params)
	if err != nil {
		return fmt.Errorf("failed to send HTTP POST request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP POST request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var responseData SendCloudResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if !responseData.Result {
		return fmt.Errorf("SMS sending failed: %s", responseData.Message)
	}

	return nil
}

// prepareParams prepares parameters for sending SMS.
func (client *SendCloudClient) prepareParams(vars map[string]string, targetPhoneNumbers []string) (url.Values, error) {
	params := url.Values{}
	params.Set("smsUser", client.SmsUser)
	params.Set("msgType", strconv.Itoa(client.MsgType))
	params.Set("phone", strings.Join(targetPhoneNumbers, ","))
	params.Set("templateId", strconv.Itoa(client.TemplateId))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))

	if len(vars) > 0 {
		varsJSON, err := json.Marshal(vars)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vars: %v", err)
		}
		params.Set("vars", string(varsJSON))
	}

	return params, nil
}

// validateConfig validates the SMS sending configuration.
func validateConfig(config SendCloudConfig) error {
	switch {
	case config.CharSet == "":
		return errors.New("charSet cannot be empty")
	case config.Server == "":
		return errors.New("server cannot be empty")
	case config.SendSMSAPI == "":
		return errors.New("sendSMSAPI cannot be empty")
	case config.MaxReceivers <= 0:
		return errors.New("maxReceivers must be greater than zero")
	}

	return nil
}

// validateSendCloudSms validates the SendCloudSms data.
func (client *SendCloudClient) validateSendCloudSms(targetPhoneNumbers []string) error {
	switch {
	case client.TemplateId == 0:
		return errors.New("templateId cannot be zero")
	case client.MsgType < 0:
		return errors.New("msgType cannot be negative")
	case len(targetPhoneNumbers) == 0:
		return errors.New("phone cannot be empty")
	}
	return nil
}

// calculateSignature calculates the signature for the request.
func calculateSignature(params url.Values, key string) string {
	sortedParams := params.Encode()
	signStr := sortedParams + "&key=" + key
	hasher := sha256.New()
	hasher.Write([]byte(signStr))
	return hex.EncodeToString(hasher.Sum(nil))
}
