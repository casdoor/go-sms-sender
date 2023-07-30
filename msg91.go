package go_sms_sender

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Msg91Client struct {
	authKey    string
	senderId   string
	templateId string
}

func GetMsg91Client(authKey string, senderId string, templateId string) (*Msg91Client, error) {
	msg91Client := &Msg91Client{
		senderId:   senderId,
		templateId: templateId,
	}

	return msg91Client, nil
}

func (m *Msg91Client) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	if len(targetPhoneNumber) < 1 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	url := "https://control.msg91.com/api/v5/flow/"

	for i := 1; i < len(targetPhoneNumber); i++ {
		if strings.HasPrefix(targetPhoneNumber[i], "+") {
			targetPhoneNumber = targetPhoneNumber[1:]
		}

		payload, err := buildPayload(m.templateId, m.senderId, "0", targetPhoneNumber[i], param)
		if err != nil {
			return fmt.Errorf("SMS build payload failed: %v", err)
		}

		err = postMsg91SendRequest(url, strings.NewReader(payload), m.authKey)
		if err != nil {
			return fmt.Errorf("send message failed：%v", err)
		}
	}

	return nil
}

func buildPayload(templateId, senderId, shortURL, mobiles string, variables map[string]string) (string, error) {
	// Create a map to hold the JSON fields
	payload := make(map[string]interface{})

	// Add the main fields
	payload["template_id"] = templateId
	payload["sender"] = senderId
	payload["short_url"] = shortURL
	payload["mobiles"] = mobiles

	// Add the variables as separate key-value pairs
	for k, v := range variables {
		payload[k] = v
	}

	// Marshal the map to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func postMsg91SendRequest(url string, payload io.Reader, authKey string) error {
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", authKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
}
