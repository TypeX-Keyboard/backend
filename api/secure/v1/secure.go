package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AcquireReq struct {
	g.Meta `path:"/acquire" tags:"安全" method:"get" summary:"申请 RSA 公钥"`
}

type AcquireRes struct {
	g.Meta    `mime:"application/json"`
	PublicKey string `json:"public_key" dc:"公钥"`
}

type SubmitReq struct {
	g.Meta `path:"/submit" tags:"安全" method:"post" summary:"上报 AES 加密密钥"`
	AesKey string `json:"aes_key" dc:"AES 加密密钥"`
}

type SubmitRes struct {
	g.Meta `mime:"application/json"`
}

type TestGetReq struct {
	g.Meta  `path:"/testGet" tags:"安全" method:"get" summary:"测试get签名"`
	Keyword string `json:"keyword"`
}

type TestGetRes struct {
	g.Meta  `mime:"application/json"`
	Keyword string `json:"keyword"`
}

type TestPostReq struct {
	g.Meta  `path:"/testPost" tags:"安全" method:"post" summary:"测试post签名"`
	Keyword string `json:"keyword"`
}

type TestPostRes struct {
	g.Meta  `mime:"application/json"`
	Keyword string `json:"keyword"`
}

type GetKeyReq struct {
	g.Meta `path:"/getKey" tags:"安全" method:"post" summary:"获取客户端key"`
	From   string `json:"from"`
}

type GetKeyRes struct {
	g.Meta  `mime:"application/json"`
	Key     string `json:"key"`
	SignKey string `json:"sign_key"`
	Expire  int64  `json:"expire"`
}

type GenKeyReq struct {
	g.Meta `path:"/genKey" tags:"安全" method:"post" summary:"重新生成key"`
}

type GenKeyRes struct {
}
