// Copyright 2022 The Casdoor Authors. All Rights Reserved.
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
// console:  https://dev.quanmwl.com/console

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
  "gitee.com/chengdu-quanming-network/quanmsms-go"
)

type QuanmSMSClient struct {
	openid     string
	apikey     string
	sign       string
	templateid string
}

func GetQuanmSMSClient(openid string, apikey string, sign string, templateid string, other []string) (*QuanmSMSClient, error) {
	return &QuanmSMSClient{
		openid: openid,
		apikey:   apikey,
		templateid: templateid,
	}, nil
}

func (c *QuanmSMSClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: code")
	}

	if len(targetPhoneNumber) == 0 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	smsSDK = quanmsms.NewSmsSDK(c.openid, c.apikey, "https")
  template_args := map[string]interface{}{"code": code}
	for _, tel := range targetPhoneNumber {
		resp, err := smsSDK.Send(tel, c.templatid, template_args)
  	if err != nil {
  		return fmt.Errorf("other error:", err.Error())
  	}
		switch string(resp.State) {
		case "201":
      // 表单信息有误或触发限发机制
			return fmt.Errorf("Refusal to send due to security reasons")
		case "202":
      // 账户重复
			return fmt.Errorf("account repeat")
		case "203":
      // 服务器异常或限流
			return fmt.Errorf("server error,Please try again later")
		case "205":
      // 请求不安全
			return fmt.Errorf("Illegal request")
		case "207":
      // 配额不足
			return fmt.Errorf("Insufficient balance")
		case "208":
      // 验签失败
			return fmt.Errorf("Verification failed")
		}
		case "209":
      // 账户被禁用接口
			return fmt.Errorf("Insufficient permissions")
		}
		case "210":
      // 账户被冻结
			return fmt.Errorf("Account frozen")
		}
		case "211":
      // 请求参数超过上限
			return fmt.Errorf("Parameter too long")
		}
		case "212":
      // 权限不足或使用了他人模板
			return fmt.Errorf("Insufficient permissions or using someone else's template")
		}
		case "213":
      // 调用状态异常
			return fmt.Errorf("status error")
		}
		case "215":
      // 内容受限
			return fmt.Errorf("Content restricted")
		}
		case "216":
      // 内容违法
			return fmt.Errorf("Content violation")
		}
	}

	return nil
}
