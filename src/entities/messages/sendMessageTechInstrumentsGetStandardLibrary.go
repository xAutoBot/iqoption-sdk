package messages

import "encoding/json"

type techInstrumentsGetStandardLibraryBody struct {
	Version        int64 `json:"version"`
	RuntimeVersion int   `json:"runtime_version"`
}

type techInstrumentsGetStandardLibrary struct {
	Name    string                                `json:"name"`
	Version string                                `json:"version"`
	Body    techInstrumentsGetStandardLibraryBody `json:"body"`
}

type sendMessageTechInstrumentsGetStandardLibrary struct {
	Name      string `json:"name"`
	RequestID string `json:"request_id"`
	LocalTime int    `json:"local_time"`
	Msg       techInstrumentsGetStandardLibrary
}

func NewSendMessageTechInstrumentsGetStandardLibrary() sendMessageTechInstrumentsGetStandardLibrary {
	body := techInstrumentsGetStandardLibraryBody{
		Version:        4657112160311,
		RuntimeVersion: 109,
	}

	msg := techInstrumentsGetStandardLibrary{
		Name:    "tech-instruments.get-standard-library",
		Version: "3.0",
		Body:    body,
	}

	return sendMessageTechInstrumentsGetStandardLibrary{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       msg,
	}
}

func (s sendMessageTechInstrumentsGetStandardLibrary) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, err
}
