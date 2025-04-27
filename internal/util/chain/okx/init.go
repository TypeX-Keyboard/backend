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
	AccessKey  string // String-type API key (follow this guide to generate an API key)
	SecretKey  string // String type secret key (follow this guide to generate a secret key)
	Passphrase string // The Passphrase you specified when you created the API key
	ProjectID  string // The project ID of your project (which can be found under Project Details)
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
