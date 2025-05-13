package ws

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 当前输入对象
type request struct {
	Event     string `json:"e"` //The name of the event
	Data      g.Map  `json:"d"` //data
	AuthToken string `json:"auth"`
}

// WResponse 输出对象
type WResponse struct {
	Event string      `json:"e"` //The name of the event
	Data  interface{} `json:"d"` //data
}

type PriceResponse struct {
	Event string              `json:"e"` //The name of the event
	Data  map[string]PoolInfo `json:"d"` //data
}

type PoolInfo struct {
	Price     float64 `json:"p"`
	Sh5m      float64 `json:"sh_5m"`
	Sh1h      float64 `json:"sh_1h"`
	Sh6h      float64 `json:"sh_6h"`
	Sh24h     float64 `json:"sh_24h"`
	Mcap      float64 `json:"mcap"`
	Fdv       float64 `json:"fdv"`
	Liquidity float64 `json:"liquidity"`
	Volume24H float64 `json:"volume_24h"`
}

type TagWResponse struct {
	Tag       string
	WResponse *WResponse
}

type UserWResponse struct {
	UserID    uint64
	WResponse *WResponse
}

type ClientWResponse struct {
	ID        string
	WResponse *WResponse
}
