package go_sms_sender

import (
	"testing"
)

func TestMockerSendMessage(t *testing.T) {
	client, err := NewSmsClient(MockSms, "", "", "", "")
	if err != nil {
		t.Fatalf("Failed to create mock client: %s", err)
	}

	err = client.SendMessage(map[string]string{"key": "value"}, "1234567890")
	if err != nil {
		t.Fatalf("Failed to send message with mock client: %s", err)
	}
}

func TestNetgsmSendMessage(t *testing.T) {
	provider := Netgsm
	accessId := ""  // KullaniciAdi
	accessKey := "" // Sifre
	sign := ""      // Baslik
	template := ""

	client, err := NewSmsClient(provider, accessId, accessKey, sign, template)
	if err != nil {
		t.Fatalf("Failed to create Netgsm client: %s", err.Error())
	}

	err = client.SendMessage(map[string]string{"key": "value"}, "5446459333")
	if err != nil {
		t.Fatalf("Failed to send message with Netgsm client: %s", err.Error())
	}
}
