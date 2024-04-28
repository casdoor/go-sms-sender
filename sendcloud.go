package go_sms_sender

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ResponseData struct {
	Result     bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Info       string `json:"info"`
}

type Config struct {
	CharSet      string
	Server       string
	SendSMSAPI   string
	SmsUser      string
	SmsKey       string
	MaxReceivers int
}

func NewConfig(smsUser, smsKey string) (Config, error) {
	if smsUser == "" {
		return Config{}, errors.New("smsUser cannot be empty")
	}
	if smsKey == "" {
		return Config{}, errors.New("smsKey cannot be empty")
	}

	return Config{
		CharSet:      "utf-8",
		Server:       "https://api.sendcloud.net",
		SendSMSAPI:   "https://api.sendcloud.net/smsapi/send",
		MaxReceivers: 100,
		SmsUser:      smsUser,
		SmsKey:       smsKey,
	}, nil
}

type SendCloudSms struct {
	TemplateId int
	MsgType    int
	Phone      []string
	Vars       map[string]string
}

func NewSendCloudSms(templateId, msgType int, phone []string, vars map[string]string) (SendCloudSms, error) {
	if templateId == 0 {
		return SendCloudSms{}, errors.New("templateId cannot be zero")
	}
	if msgType < 0 {
		return SendCloudSms{}, errors.New("msgType cannot be negative")
	}
	if len(phone) == 0 {
		return SendCloudSms{}, errors.New("phone cannot be empty")
	}
	return SendCloudSms{
		TemplateId: templateId,
		MsgType:    msgType,
		Phone:      phone,
		Vars:       vars,
	}, nil
}

func SendMessage(sms SendCloudSms, config Config) error {

	if err := ValidateSendCloudSms(sms); err != nil {
		return err
	}
	if err := ValidateConfig(config); err != nil {
		return err
	}

	params := url.Values{}
	params.Set("smsUser", config.SmsUser)
	params.Set("msgType", strconv.Itoa(sms.MsgType))
	params.Set("phone", strings.Join(sms.Phone, ","))
	params.Set("templateId", strconv.Itoa(sms.TemplateId))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))

	if len(sms.Vars) > 0 {
		varsJSON, err := json.Marshal(sms.Vars)
		if err != nil {
			return fmt.Errorf("failed to marshal vars: %v", err)
		}
		params.Set("vars", string(varsJSON))
	}

	signature := calculateSignature(params, config.SmsKey)
	params.Set("signature", signature)

	resp, err := http.PostForm(config.SendSMSAPI, params)
	if err != nil {
		return fmt.Errorf("failed to send HTTP POST request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return nil
}

func ValidateConfig(config Config) error {
	if config.CharSet == "" {
		return errors.New("charSet cannot be empty")
	}
	if config.Server == "" {
		return errors.New("server cannot be empty")
	}
	if config.SendSMSAPI == "" {
		return errors.New("sendSMSAPI cannot be empty")
	}
	if config.SmsUser == "" {
		return errors.New("smsUser cannot be empty")
	}
	if config.SmsKey == "" {
		return errors.New("smsKey cannot be empty")
	}
	if config.MaxReceivers <= 0 {
		return errors.New("maxReceivers must be greater than zero")
	}
	return nil
}

func ValidateSendCloudSms(sms SendCloudSms) error {
	if sms.TemplateId == 0 {
		return errors.New("templateId cannot be zero")
	}
	if sms.MsgType < 0 {
		return errors.New("msgType cannot be negative")
	}
	if len(sms.Phone) == 0 {
		return errors.New("phone cannot be empty")
	}
	return nil
}

func calculateSignature(params url.Values, key string) string {

	sortedParams := params.Encode()
	signStr := sortedParams + key
	hasher := md5.New()
	hasher.Write([]byte(signStr))
	signature := hex.EncodeToString(hasher.Sum(nil))
	return signature
}
