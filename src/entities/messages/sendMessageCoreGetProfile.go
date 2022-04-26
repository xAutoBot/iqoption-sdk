package messages

type coreGetProfileBody struct {
}

type coreGetProfileMsg struct {
	Name    string             `json:"name"`
	Version string             `json:"version"`
	Body    coreGetProfileBody `json:"body"`
}
type sendMessageCoreGetProfile struct {
	Name      string            `json:"name"`
	RequestID string            `json:"request_id"`
	LocalTime int               `json:"local_time"`
	Msg       coreGetProfileMsg `json:"msg"`
}

func NewSendMessageCoreGetProfile() sendMessageCoreGetProfile {
	coreGetProfileBody := coreGetProfileBody{}

	coreGetProfileMsg := coreGetProfileMsg{
		Name:    "core.get-profile",
		Version: "1.0",
		Body:    coreGetProfileBody,
	}

	return sendMessageCoreGetProfile{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       coreGetProfileMsg,
	}
}
