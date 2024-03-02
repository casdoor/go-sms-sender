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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OsonClient struct {
	Endpoint         string
	SenderID         string
	SecretAccessHash string
	Sign             string
	Message          string
}

func GetOsonClient(senderId, secretAccessHash, sign, message string) (*OsonClient, error) {
	return &OsonClient{
		Endpoint:         "https://api.osonsms.com/sendsms_v1.php",
		SenderID:         senderId,
		SecretAccessHash: secretAccessHash,
		Sign:             sign,
		Message:          message,
	}, nil
}

func (c *OsonClient) SendMessage(param map[string]string, targetPhoneNumber ...string) (err error) {
	// Init http client for make request to sms center.
	// Set a timeout of 25+ seconds to ensure that the
	// response from the SMS center has been processed.
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	txnID := uuid.New()
	strHash := strings.Join([]string{txnID.String(), c.SenderID, c.Sign, targetPhoneNumber[0], c.SecretAccessHash}, ";")

	urlLink, err := url.Parse(c.Endpoint)
	if err != nil {
		return
	}

	urlParams := url.Values{}
	urlParams.Add("from", c.Sign)
	urlParams.Add("phone_number", targetPhoneNumber[0])
	urlParams.Add("msg", c.Message)
	urlParams.Add("str_hash", strHash)
	urlParams.Add("txn_id", txnID.String())
	urlParams.Add("login", c.SenderID)

	urlLink.RawQuery = urlParams.Encode()

	request, err := http.NewRequest(http.MethodGet, urlLink.String(), nil)
	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sms service returned error status not 200: Status Code: %d Error:%v", resp.StatusCode, fmt.Sprint(result[:]))
	}

	return
}
