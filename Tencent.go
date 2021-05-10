package go_sms_sender

import (
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

type TencentClient struct {
	core     *sms.Client
	appId    string
	sign     string
	template string
}

func GetTencentClient(accessId, accessKey, sign, region, templateId string, appId []string) *TencentClient {
	if len(appId) < 1 {
		panic(SmsError{"Tencent SMS Client: missing parameter appId"})
	}

	credential := common.NewCredential(accessId, accessKey)
	config := profile.NewClientProfile()
	config.HttpProfile.ReqMethod = "POST"

	client, err := sms.NewClient(credential, region, config)
	if err != nil {
		panic(err)
	}
	return &TencentClient{
		core:     client,
		appId:    appId[0],
		sign:     sign,
		template: templateId,
	}
}

func (c *TencentClient) SendMessage(param map[string]string, targetPhoneNumber ...string) {
	var paramArray []string
	index := 0
	for {
		value := param[strconv.Itoa(index)]
		if len(value) == 0 {
			break
		}
		paramArray = append(paramArray, value)
		index++
	}

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppid = common.StringPtr(c.appId)
	request.Sign = common.StringPtr(c.sign)
	request.TemplateParamSet = common.StringPtrs(paramArray)
	request.TemplateID = common.StringPtr(c.template)
	request.PhoneNumberSet = common.StringPtrs(targetPhoneNumber)
	_, err := c.core.SendSms(request)
	if err != nil {
		panic(err)
	}
}
