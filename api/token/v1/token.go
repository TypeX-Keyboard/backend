package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/model"
)

type GetTokenDetailReq struct {
	g.Meta  `path:"/get-token-detail" tags:"Tokens" method:"get" summary:"Get token details"`
	Address string `json:"address" dc:"Token address" required:"true"`
}

type GetTokenDetailRes model.TokenDetail

type CheckAddressReq struct {
	g.Meta  `path:"/check-address" tags:"Tokens" method:"get" summary:"Check the address"`
	Address string `json:"address" dc:"address" required:"true"`
}

type CheckAddressRes model.CheckAddress
