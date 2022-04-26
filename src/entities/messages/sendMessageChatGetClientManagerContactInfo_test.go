package messages

import "testing"

func TestNewSendMessageChatGetClientManagerContactInfo(t *testing.T) {

	jsonExpected := `{"name":"sendMessage","request_id":"4","local_time":4199,"msg":{"name":"chat.get-client-manager-contact-info","version":"1.0"}}`
	sendMessageChatGetClientManagerContactInfo := NewSendMessageChatGetClientManagerContactInfo()

	if j, _ := sendMessageChatGetClientManagerContactInfo.Json(); string(j) != jsonExpected {
		t.Errorf("json is not equal")
	}
}
