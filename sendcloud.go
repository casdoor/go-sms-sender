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

// ResponseData structure holds the response from SendCloud API.
type ResponseData struct {
	Result     bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Info       string `json:"info"`
}

// Config holds the configuration for the SMS sending service.
type Config struct {
	CharSet      string
	Server       string
	SendSMSAPI   string
	SmsUser      string
	SmsKey       string
	MaxReceivers int
}

// NewConfig creates a new configuration with the provided user and key.
func NewConfig(smsUser, smsKey string) (*Config, error) {
	if smsUser == "" {
		return nil, errors.New("smsUser cannot be empty")
	}
	if smsKey == "" {
		return nil, errors.New("smsKey cannot be empty")
	}

	return &Config{
		CharSet:      "utf-8",
		Server:       "https://api.sendcloud.net",
		SendSMSAPI:   "/smsapi/send",
		MaxReceivers: 100,
		SmsUser:      smsUser,
		SmsKey:       smsKey,
	}, nil
}

// SendCloudSms represents the data required to send an SMS message.
type SendCloudSms struct {
	TemplateId int
	MsgType    int
	Phone      []string
	Vars       map[string]string
}

// NewSendCloudSms creates a new SendCloudSms instance with the provided parameters.
func NewSendCloudSms(templateId, msgType int, phone []string, vars map[string]string) (*SendCloudSms, error) {
	if templateId <= 0 {
		return nil, errors.New("templateId must be greater than zero")
	}
	if msgType < 0 {
		return nil, errors.New("msgType cannot be negative")
	}
	if len(phone) == 0 {
		return nil, errors.New("phone cannot be empty")
	}
	return &SendCloudSms{
		TemplateId: templateId,
		MsgType:    msgType,
		Phone:      phone,
		Vars:       vars,
	}, nil
}

// SendMessage sends an SMS using SendCloud API.
func SendMessage(sms SendCloudSms, config Config) error {
	if err := validateSendCloudSms(sms); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	params, err := prepareParams(sms, config)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	signature := calculateSignature(params, config.SmsKey)
	params.Set("signature", signature)

	resp, err := http.PostForm(config.SendSMSAPI, params)
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

	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if !responseData.Result {
		return fmt.Errorf("SMS sending failed: %s", responseData.Message)
	}

	return nil
}

// prepareParams prepares parameters for sending SMS.
func prepareParams(sms SendCloudSms, config Config) (url.Values, error) {
	params := url.Values{}
	params.Set("smsUser", config.SmsUser)
	params.Set("msgType", strconv.Itoa(sms.MsgType))
	params.Set("phone", strings.Join(sms.Phone, ","))
	params.Set("templateId", strconv.Itoa(sms.TemplateId))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))

	if len(sms.Vars) > 0 {
		varsJSON, err := json.Marshal(sms.Vars)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vars: %v", err)
		}
		params.Set("vars", string(varsJSON))
	}

	return params, nil
}

// validateConfig validates the SMS sending configuration.
func validateConfig(config Config) error {
	switch {
	case config.CharSet == "":
		return errors.New("charSet cannot be empty")
	case config.Server == "":
		return errors.New("server cannot be empty")
	case config.SendSMSAPI == "":
		return errors.New("sendSMSAPI cannot be empty")
	case config.SmsUser == "":
		return errors.New("smsUser cannot be empty")
	case config.SmsKey == "":
		return errors.New("smsKey cannot be empty")
	case config.MaxReceivers <= 0:
		return errors.New("maxReceivers must be greater than zero")
	}

	return nil
}

// validateSendCloudSms validates the SendCloudSms data.
func validateSendCloudSms(sms SendCloudSms) error {
	switch {
	case sms.TemplateId == 0:
		return errors.New("templateId cannot be zero")
	case sms.MsgType < 0:
		return errors.New("msgType cannot be negative")
	case len(sms.Phone) == 0:
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
