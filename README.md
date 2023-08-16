# go-sms-sender

[![Go Report Card](https://goreportcard.com/badge/github.com/casdoor/go-sms-sender)](https://goreportcard.com/report/github.com/casdoor/go-sms-sender)
[![Go](https://github.com/casdoor/go-sms-sender/actions/workflows/ci.yml/badge.svg)](https://github.com/casdoor/go-sms-sender/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/casdoor/go-sms-sender.svg)](https://pkg.go.dev/github.com/casdoor/go-sms-sender)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/casdoor/go-sms-sender)

This is a powerful open-source library for sending SMS message, which will help you to easily integrate with the popular SMS providers. And it has been applied to [Casdoor](https://github.com/casdoor/casdoor), if you still don’t know how to use it after reading README.md, you can refer to it.

We support the following SMS providers, welcome to contribute.

- [Twilio](https://www.twilio.com)
- [Amazon SNS](https://aws.amazon.com/sns/)
- [Azure ACS](https://azure.microsoft.com/en-us/products/communication-services)
- [GCCPAY](https://gccpay.com/)
- [Infobip](https://www.infobip.com/)
- [SUBMAIL](https://en.mysubmail.com/)
- [SmsBao](https://www.smsbao.com/)
- [Alibaba Cloud](https://www.aliyun.com/product/sms)
- [Tencent Cloud](https://cloud.tencent.com/document/product/382)
- [Baidu Cloud](https://cloud.baidu.com/product/sms.html)
- [VolcEngine](https://www.volcengine.com/product/cloud-sms)
- [Huawei Cloud](https://www.huaweicloud.com/product/msgsms.html)
- [UCloud](https://www.ucloud.cn/site/product/usms.html)
- [Huyi](https://www.ihuyi.com/)
- [Netgsm](https://www.netgsm.com.tr/) 🇹🇷 🇹🇷 🇹🇷

## Installation

Use `go get` to install：

```
go get github.com/casdoor/go-sms-sender
```

## How to use

### Create Client

Different SMS providers need to provide different configuration, but we support a unit API as below to init and get the SMS client.

```go
func NewSmsClient(provider string, accessId string, accessKey string, sign string, template string, other ...string) (SmsClient, error)
```

- `provider` the name of SMS provider, such as `Aliyun SMS`
- `accessId`
- `accessKey`
- `sign` the sign name
- `template` the template code
- `other` other configuration

### Send Message

After initializing the SMS client, we can use the following API to send message.

```go
SendMessage(param map[string]string, targetPhoneNumber ...string) error
```

- `param` the parameters in the SMS template, such as 6 random numbers
- `targetPhoneNumber` the receivers, such as `+8612345678910`

## Example

### Twilio

Please get necessary information from Twilio [console](https://console.twilio.com/)

```go
package main

import "github.com/casdoor/go-sms-sender"

func main() {
	client, err := go_sms_sender.NewSmsClient(go_sms_sender.Twilio, "ACCOUNT_SID", "AUTH_TOKEN", "", "TEMPLATE_CODE")
	if err != nil {
		panic(err)
	}

	params := map[string]string{}
	params["code"] = "123456"
	phoneNumer := "+8612345678910"
	err = client.SendMessage(params, phoneNumer)
	if err != nil {
		panic(err)
	}
}
```

### Aliyun

Before you begin, you need to sign up for an Aliyun account and retrieve your [Credentials](https://usercenter.console.aliyun.com/#/manage/ak).

```go
package main

import "github.com/casdoor/go-sms-sender"

func main() {
	client, err := go_sms_sender.NewSmsClient(go_sms_sender.Aliyun, "ACCESS_KEY_ID", "ACCESS_KEY_SECRET", "SIGN_NAME", "TEMPLATE_CODE")
	if err != nil {
		panic(err)
	}

	params := map[string]string{}
	params["code"] = "473956"
	phoneNumer := "+8612345678910"
	err = client.SendMessage(params, phoneNumer)
	if err != nil {
		panic(err)
	}
}
```

### Tencent Cloud

```go
package main

import "github.com/casdoor/go-sms-sender"

func main() {
	client, err := go_sms_sender.NewSmsClient(go_sms_sender.TencentCloud, "secretId", "secretKey", "SIGN_NAME", "TEMPLATE_CODE", "APP_ID")
	if err != nil {
		panic(err)
	}

	params := map[string]string{}
	params["0"] = "473956"
	phoneNumer := "+8612345678910"
	err = client.SendMessage(params, phoneNumer)
	if err != nil {
		panic(err)
	}
}
```

### Netgsm

- yourAccessId: is KullaniciAdi
- yourAccessKey: is Sifre
- yourSign: is Baslik

```go
package main

import "github.com/casdoor/go-sms-sender"

func main() {
	client, err := go_sms_sender.NewSmsClient(go_sms_sender.Netgsm, "yourAccessId", "yourAccessKey", "yourSign", "yourTemplate")
	if err != nil {
		panic(err)
	}

	params := map[string]string{}
	params["param1"] = "value1"
	params["param2"] = "value2"
	phoneNumer := "+8612345678910"
	err = client.SendMessage(params, phoneNumer)
	if err != nil {
		panic(err)
	}
}
```

### Running Tests

To run tests for the `go-sms-sender` library, navigate to the root folder of the project in your terminal and execute the following command:

```sh
go test -v ./...
```

you can modify mock_test.go file to mock an other tests

## License

This project is licensed under the [Apache 2.0 license](LICENSE).
