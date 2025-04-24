package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AcquireReq struct {
	g.Meta `path:"/acquire" tags:"secure" method:"get" summary:"Request an RSA public key"`
}

type AcquireRes struct {
	g.Meta    `mime:"application/json"`
	PublicKey string `json:"public_key" dc:"PublicKey"`
}

type SubmitReq struct {
	g.Meta `path:"/submit" tags:"secure" method:"post" summary:"Report the AES encryption key"`
	AesKey string `json:"aes_key" dc:"AES encryption key"`
}

type SubmitRes struct {
	g.Meta `mime:"application/json"`
}

type TestGetReq struct {
	g.Meta  `path:"/testGet" tags:"secure" method:"get" summary:"Test the get signature"`
	Keyword string `json:"keyword"`
}

type TestGetRes struct {
	g.Meta  `mime:"application/json"`
	Keyword string `json:"keyword"`
}

type TestPostReq struct {
	g.Meta  `path:"/testPost" tags:"secure" method:"post" summary:"Test POST signatures"`
	Keyword string `json:"keyword"`
}

type TestPostRes struct {
	g.Meta  `mime:"application/json"`
	Keyword string `json:"keyword"`
}

type GetKeyReq struct {
	g.Meta `path:"/getKey" tags:"secure" method:"post" summary:"Obtain the client key"`
	From   string `json:"from"`
}

type GetKeyRes struct {
	g.Meta  `mime:"application/json"`
	Key     string `json:"key"`
	SignKey string `json:"sign_key"`
	Expire  int64  `json:"expire"`
}

type GenKeyReq struct {
	g.Meta `path:"/genKey" tags:"secure" method:"post" summary:"Regenerate the key"`
}

type GenKeyRes struct {
}
