package responseMessage

type ResponnseCandleGenerated struct {
	Name             string `json:"name"`
	MicroserviceName string `json:"microserviceName"`
	Msg              Candle `json:"msg"`
}
type Candle struct {
	ActiveID int     `json:"active_id"`
	Size     int     `json:"size"`
	At       int64   `json:"at"`
	From     int64   `json:"from"`
	To       int64   `json:"to"`
	ID       int     `json:"id"`
	Open     float64 `json:"open"`
	Close    float64 `json:"close"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Ask      float64 `json:"ask"`
	Bid      float64 `json:"bid"`
	Volume   int64   `json:"volume"`
	Phase    string  `json:"phase"`
}
