package go_sms_sender

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NetgsmClient struct {
	accessId   string
	accessKey  string
	sign       string
	template   string
	httpClient *http.Client
}

type NetgsmResponse struct {
	Code  string `xml:"main>code"`
	JobID string `xml:"main>jobID"`
	Error string `xml:"main>error"`
}

func GetNetgsmClient(accessId, accessKey, sign, template string) (*NetgsmClient, error) {
	return &NetgsmClient{
		accessId:   accessId,
		accessKey:  accessKey,
		sign:       sign,
		template:   template,
		httpClient: &http.Client{},
	}, nil
}

func (c *NetgsmClient) SendMessage(param map[string]string, targetPhoneNumbers ...string) error {
	for _, phoneNumber := range targetPhoneNumbers {
		data := fmt.Sprintf(`
<mainbody>
   <header>
       <usercode>%s</usercode>
       <password>%s</password>
       <msgheader>%s</msgheader>
   </header>
   <body>
       <msg>
           <![CDATA[%s]]>
       </msg>
       <no>%s</no>
   </body>
</mainbody>`, c.accessId, c.accessKey, c.sign, c.template, phoneNumber)

		headers := map[string]string{
			"Content-Type": "application/xml",
		}

		respBody, err := c.postXML("https://api.netgsm.com.tr/sms/send/otp", data, headers)
		if err != nil {
			return err
		}

		var netgsmResponse NetgsmResponse
		if err := xml.Unmarshal([]byte(respBody), &netgsmResponse); err != nil {
			return err
		}

		if netgsmResponse.Code != "0" {
			return errors.New(netgsmResponse.Error)
		}
	}
	return nil
}

func (c *NetgsmClient) postXML(url, xmlData string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
