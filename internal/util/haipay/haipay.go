package haipay

type sHaiPay struct {
}

type IHaiPay interface {
}

func init() {
	RegisterHaiPay(New())
}

func New() IHaiPay {
	return &sHaiPay{}
}

var (
	localHaiPay IHaiPay
)

func HaiPay() IHaiPay {
	if localHaiPay == nil {
		panic("implement not found for interface IHaiPay, forgot register?")
	}
	return localHaiPay
}

func RegisterHaiPay(i IHaiPay) {
	localHaiPay = i
}
