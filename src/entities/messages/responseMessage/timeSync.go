package responseMessage

type TimeSync struct {
	Name string `json:"name"`
	Msg  int64  `json:"msg"`
}

//{"name":"timeSync","msg":1650576081919}
