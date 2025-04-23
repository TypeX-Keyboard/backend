package okx

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/net/context"
)

type IOkx interface {
	CreateWalletAccount(ctx context.Context, accountReq CreateWalletReq) (res CreateWalletRes, err error)
	GetWalletAccount(ctx context.Context, limit int, cursor int) (res GetWalletAccountRes, err error)
	GetPrice(ctx context.Context, input []GetPriceInput) (*PriceRes, error)
	GetSupportedChains(ctx context.Context) ([]Chain, error)
	TransactionsByAddress(ctx context.Context, chain, address, tokenAddress, begin, end, cursor, limit string) (*TxHistoryRes, error)
	TransactionDetailByTxHash(ctx context.Context, chain, txHash string) (*TxHashRes, error)
	TotalValueByAddress(ctx context.Context, chain, address, assetType string, excludeRiskToken ...bool) (*TotalValueByAddressRes, error)
	AllTokenBalancesByAddress(ctx context.Context, chain, address, filter string) (*AllTokenBalancesByAddressRes, error)
	TokenDetail(ctx context.Context, chainIndex, tokenAddress string) (*TokenDetailRes, error)
	TokenList(ctx context.Context, chainId string) (*TokenListRes, error)
	TokenBalancesByAddress(ctx context.Context, chain, address, tokenAddresses string) (*AllTokenBalancesByAddressRes, error)
}

type sOkx struct {
	okxConf []okxConf
}

type okxConf struct {
	AccessKey  string // 字符串类型的 API key (遵循这个指南来生成一个 API key)
	SecretKey  string // 字符串类型的 Secret key (遵循这个指南来生成一个 Secret key)
	Passphrase string // 你在创建 API key 时指定的 Passphrase
	ProjectID  string // 你的项目的项目 ID (可在项目详细信息下找到)
}

var localOkx IOkx

func Okx() IOkx {
	if localOkx == nil {
		panic("implement not found for interface IOkx, forgot register?")
	}
	return localOkx
}

func RegisterOkx(i IOkx) {
	localOkx = i
}

func init() {
	RegisterOkx(New())
}

func New() IOkx {
	var conf []okxConf
	c := g.Config().MustGet(gctx.GetInitCtx(), "okxs")
	c.Structs(&conf)
	return &sOkx{
		okxConf: conf,
	}
}
