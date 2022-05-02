package responseMessage

type ClientPriceGenerated struct {
	Name             string                  `json:"name"`
	MicroserviceName string                  `json:"microserviceName"`
	Msg              ClientPriceGeneratedMsg `json:"msg"`
}
type call struct {
	Symbol string  `json:"symbol"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
}
type put struct {
	Symbol string  `json:"symbol"`
	Bid    float64 `json:"bid"`
	Ask    float64 `json:"ask"`
}
type prices struct {
	Strike string `json:"strike"`
	Put    put    `json:"put"`
	Call   call   `json:"call"`
}
type ClientPriceGeneratedMsg struct {
	InstrumentIndex int      `json:"instrument_index"`
	InstrumentType  string   `json:"instrument_type"`
	AssetID         int      `json:"asset_id"`
	UserGroupID     int      `json:"user_group_id"`
	QuoteTime       string   `json:"quote_time"`
	Prices          []prices `json:"prices"`
}
