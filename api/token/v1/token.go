package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/model"
)

type GetTokenDetailReq struct {
	g.Meta  `path:"/get-token-detail" tags:"代币" method:"get" summary:"获取代币详情"`
	Address string `json:"address" dc:"token地址" required:"true"`
}

type GetTokenDetailRes model.TokenDetail

type CheckAddressReq struct {
	g.Meta  `path:"/check-address" tags:"代币" method:"get" summary:"检查地址"`
	Address string `json:"address" dc:"地址" required:"true"`
}

type CheckAddressRes model.CheckAddress
