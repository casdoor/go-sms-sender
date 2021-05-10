package go_sms_sender

type SmsClient interface {
	SendMessage(param map[string]string, targetPhoneNumber ...string)
}

func NewSmsClient(provider, accessId, accessKey, sign, region, template string, other ...string) SmsClient {
	switch provider {
	case "aliyun":
		return GetAliyunClient(accessId, accessKey, sign, region, template)
	case "tencent":
		return GetTencentClient(accessId, accessKey, sign, region, template, other)
	default:
		return nil
	}
}

type SmsError struct {
	errorText string
}

func (e SmsError) Error() string {
	return e.errorText
}
