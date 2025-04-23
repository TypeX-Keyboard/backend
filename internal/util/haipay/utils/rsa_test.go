package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/google/uuid"
	"testing"
)

var ctx = gctx.New()

func TestCollect(t *testing.T) {
	privateKey := g.Config().MustGet(ctx, "haipay.privateKey").String()

	request := map[string]interface{}{
		"appId":      g.Config().MustGet(ctx, "haipay.appId").String(),
		"orderId":    uuid.New().ID(),
		"amount":     "300",
		"phone":      "09230219312",
		"email":      "23423@qq.com",
		"name":       "test",
		"inBankCode": "PH_QRPH_DYNAMIC",
		"payType":    "QR",
	}

	content := GetSign(request, g.Config().MustGet(ctx, "haipay.secretKey").String())
	signature, err := SignWithPrivateKey(content, privateKey)
	if err != nil {
		panic(err)
	}
	request["sign"] = signature
	fmt.Println("签名:", signature)
	g.Dump(request)
	c := g.Client()
	c.SetHeader("Content-Type", "application/json")
	post, err := c.Post(ctx, "https://uat-interface.haipay.asia/php/collect/apply", request)
	if err != nil {
		panic(err)
	}
	defer post.Body.Close()
	g.Dump(post.ReadAllString())
}

func TestPay(t *testing.T) {
	privateKey := g.Config().MustGet(ctx, "haipay.privateKey").String()

	request := map[string]interface{}{
		"appId":       g.Config().MustGet(ctx, "haipay.appId").String(),
		"orderId":     uuid.New().ID(),
		"accountType": "EWALLET",
		"amount":      "300",
		"phone":       "09230219312",
		"email":       "23423@qq.com",
		"name":        "test",
		"bankCode":    "PH_QRPH_DYNAMIC",
		"payType":     "QR",
	}

	content := GetSign(request, g.Config().MustGet(ctx, "haipay.secretKey").String())
	signature, err := SignWithPrivateKey(content, privateKey)
	if err != nil {
		panic(err)
	}
	request["sign"] = signature
	fmt.Println("签名:", signature)
	g.Dump(request)
	c := g.Client()
	c.SetHeader("Content-Type", "application/json")
	post, err := c.Post(ctx, "https://uat-interface.haipay.asia/php/pay/apply", request)
	if err != nil {
		panic(err)
	}
	defer post.Body.Close()
	g.Dump(post.ReadAllString())
}
