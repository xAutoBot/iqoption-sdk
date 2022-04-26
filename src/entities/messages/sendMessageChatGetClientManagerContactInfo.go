package messages

import "encoding/json"

type chatGetClientManagerContactInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type sendMessageChatGetClientManagerContactInfo struct {
	Name      string                          `json:"name"`
	RequestID string                          `json:"request_id"`
	LocalTime int                             `json:"local_time"`
	Msg       chatGetClientManagerContactInfo `json:"msg"`
}

func NewSendMessageChatGetClientManagerContactInfo() sendMessageChatGetClientManagerContactInfo {

	chatGetClientManagerContactInfo := chatGetClientManagerContactInfo{
		Name:    "chat.get-client-manager-contact-info",
		Version: "1.0",
	}
	return sendMessageChatGetClientManagerContactInfo{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       chatGetClientManagerContactInfo,
	}
}

func (s sendMessageChatGetClientManagerContactInfo) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return j, nil
}
