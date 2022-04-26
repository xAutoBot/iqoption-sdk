package messages

import "encoding/json"

type sendMessageCandleGenerate struct {
	Name      string            `json:"name"`
	RequestID string            `json:"request_id"`
	LocalTime int               `json:"local_time"`
	Msg       candleGenerateMsg `json:"msg"`
}

type candleGenerateMsg struct {
	Name   string                  `json:"name"`
	Params candleGenerateMsgParams `json:"params"`
}

type candleGenerateMsgParams struct {
	RoutingFilters routingFilters `json:"routingFilters"`
}
type routingFilters struct {
	ActiveID int `json:"active_id"`
	Size     int `json:"size"`
}

func NewSendMessageStartCandleGenerate(activeID, candleSize int) sendMessageCandleGenerate {
	routingFilters := routingFilters{
		ActiveID: activeID,
		Size:     candleSize,
	}
	candleGenerateMsgParams := candleGenerateMsgParams{
		RoutingFilters: routingFilters,
	}
	candleGenerateMsg := candleGenerateMsg{
		Name:   "candle-generated",
		Params: candleGenerateMsgParams,
	}

	return sendMessageCandleGenerate{
		Name:      "subscribeMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       candleGenerateMsg,
	}
}
func NewSendMessageStopCandleGenerate(activeID, candleSize int) sendMessageCandleGenerate {

	routingFilters := routingFilters{
		ActiveID: activeID,
		Size:     candleSize,
	}

	candleGenerateMsgParams := candleGenerateMsgParams{
		RoutingFilters: routingFilters,
	}

	candleGenerateMsg := candleGenerateMsg{
		Name:   "candle-generated",
		Params: candleGenerateMsgParams,
	}

	return sendMessageCandleGenerate{
		Name:      "unsubscribeMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       candleGenerateMsg,
	}
}

func (s sendMessageCandleGenerate) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
