package messages

import "encoding/json"

type sendMessagePriceSplitterClientPriceGenerated struct {
	Name      string                               `json:"name"`
	RequestID string                               `json:"request_id"`
	LocalTime int                                  `json:"local_time"`
	Msg       priceSplitterClientPriceGeneratedMsg `json:"msg"`
}
type priceSplitterClientPriceGeneratedRoutingFilters struct {
	InstrumentType  string `json:"instrument_type"`
	AssetID         int    `json:"asset_id"`
	InstrumentIndex int    `json:"-"`
}
type priceSplitterClientPriceGeneratedParams struct {
	RoutingFilters priceSplitterClientPriceGeneratedRoutingFilters `json:"routingFilters"`
}
type priceSplitterClientPriceGeneratedMsg struct {
	Name    string                                  `json:"name"`
	Version string                                  `json:"version"`
	Params  priceSplitterClientPriceGeneratedParams `json:"params"`
}

func NewSendMessageStartPriceSplitterClientPriceGenerated(activeID int) sendMessagePriceSplitterClientPriceGenerated {
	priceSplitterClientPriceGeneratedRoutingFilters := priceSplitterClientPriceGeneratedRoutingFilters{
		InstrumentType:  "digital-option",
		AssetID:         activeID,
		InstrumentIndex: 1234,
	}
	priceSplitterClientPriceGeneratedParams := priceSplitterClientPriceGeneratedParams{
		RoutingFilters: priceSplitterClientPriceGeneratedRoutingFilters,
	}
	priceSplitterClientPriceGeneratedMsg := priceSplitterClientPriceGeneratedMsg{
		Name:    "price-splitter.client-price-generated",
		Version: "1.0",
		Params:  priceSplitterClientPriceGeneratedParams,
	}

	return sendMessagePriceSplitterClientPriceGenerated{
		Name:      "subscribeMessage",
		RequestID: GetRequestId(),
		LocalTime: 13233637,
		Msg:       priceSplitterClientPriceGeneratedMsg,
	}
}
func NewSendMessageStopPriceSplitterClientPriceGenerated(activeID int) sendMessagePriceSplitterClientPriceGenerated {
	priceSplitterClientPriceGeneratedRoutingFilters := priceSplitterClientPriceGeneratedRoutingFilters{
		InstrumentType:  "digital-option",
		AssetID:         activeID,
		InstrumentIndex: 1234,
	}
	priceSplitterClientPriceGeneratedParams := priceSplitterClientPriceGeneratedParams{
		RoutingFilters: priceSplitterClientPriceGeneratedRoutingFilters,
	}
	priceSplitterClientPriceGeneratedMsg := priceSplitterClientPriceGeneratedMsg{
		Name:    "price-splitter.client-price-generated",
		Version: "1.0",
		Params:  priceSplitterClientPriceGeneratedParams,
	}

	return sendMessagePriceSplitterClientPriceGenerated{
		Name:      "unsubscribeMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       priceSplitterClientPriceGeneratedMsg,
	}
}

func (s sendMessagePriceSplitterClientPriceGenerated) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
