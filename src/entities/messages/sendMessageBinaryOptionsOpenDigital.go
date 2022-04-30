package messages

import (
	"encoding/json"
	"fmt"
)

type sendMessageBinaryOptionsOpenDigital struct {
	Name      string                      `json:"name"`
	RequestID string                      `json:"request_id"`
	LocalTime int                         `json:"local_time"`
	Msg       binaryOptionsOpenDigitalMsg `json:"msg"`
}

type binaryOptionsOpenDigitalMsg struct {
	Name    string                       `json:"name"`
	Version string                       `json:"version"`
	Body    BinaryOptionsOpenDigitalBody `json:"body"`
}

type BinaryOptionsOpenDigitalBody struct {
	UserBalanceID   int    `json:"user_balance_id"`
	InstrumentID    string `json:"instrument_id"`
	Amount          string `json:"amount"`
	InstrumentIndex int    `json:"instrument_index"`
	AssetID         int    `json:"asset_id"`
}

func NewBinaryOptionsOpenDigitalBody(userBalanceId int, instrumentID string, amount float32, instrumentIndex int, assetID int) BinaryOptionsOpenDigitalBody {
	amountStr := fmt.Sprintf("%2f", amount)
	return BinaryOptionsOpenDigitalBody{
		UserBalanceID:   userBalanceId,
		InstrumentID:    instrumentID,
		Amount:          amountStr,
		InstrumentIndex: instrumentIndex,
		AssetID:         assetID,
	}
}

func NewSendMessageBinaryOptionsOpenDigital(b BinaryOptionsOpenDigitalBody) sendMessageBinaryOptionsOpenDigital {

	binaryOptionsOpenDigitalMsg := binaryOptionsOpenDigitalMsg{
		Name:    "digital-options.place-digital-option",
		Version: "2.0",
		Body:    b,
	}

	return sendMessageBinaryOptionsOpenDigital{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       binaryOptionsOpenDigitalMsg,
	}
}

func (s sendMessageBinaryOptionsOpenDigital) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
