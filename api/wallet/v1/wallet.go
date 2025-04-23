package v1

import (
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/util/chain/okx"

	"github.com/gogf/gf/v2/frame/g"
)

type GetSupportedChainsReq struct {
	g.Meta `path:"/supported-chains" tags:"钱包" method:"get" summary:"查询支持的链"`
}

type GetSupportedChainsRes model.SupportedChainsRes

type GetTransactionsByAddressReq struct {
	g.Meta       `path:"/transactions-by-address" tags:"钱包" method:"get" summary:"地址维度查询交易历史"`
	Address      string `json:"address" dc:"具体某条链的账户地址" required:"true"`
	Chains       string `json:"chains" dc:"筛选需要查询交易历史的链，多条链以,分隔" required:"true" d:"501"`
	TokenAddress string `json:"tokenAddress" dc:"代币地址 1：传空代表查询对应链的主链币。2：传具体的代币合约地址，代表查询对应的代币。3：不传，代表查询主链币和所有代币。" required:"false"`
	Begin        string `json:"begin" dc:"开始时间，查询晚于该时间的交易历史。Unix 时间戳，用毫秒表示" required:"false"`
	End          string `json:"end" dc:"结束时间，查询早于该时间的交易历史。若 begin 和 end 都不传，查询当前时间以前的交易历史。Unix 时间戳，用毫秒表示" required:"false"`
	Cursor       string `json:"cursor" dc:"游标" required:"false"`
	Limit        string `json:"limit" dc:"返回条数，默认返回最近的 20 条。多链查询最多支持 20 条，单链查询最多 100 条" required:"false"`
}

type GetTransactionsByAddressRes struct {
	Cursor          string                        `json:"cursor"`
	TransactionList []model.TransactionsByAddress `json:"transactionList"`
}

type GetTransactionByHashReq struct {
	g.Meta  `path:"/transaction-by-hash" tags:"钱包" method:"get" summary:"交易哈希查询交易详情"`
	Chain   string `json:"chain" dc:"筛选需要查询交易历史的链" required:"true" d:"501"`
	TxHash  string `json:"txHash" dc:"交易hash"`
	Address string `json:"address" dc:"具体某条链的账户地址" required:"true"`
}

type GetTransactionByHashRes struct {
	PlatformFee     string `json:"platformFee" dc:"平台手续费"` // 平台手续费
	Action          int    `json:"action" dc:"1=接收，2=发送，3=买入，4=卖出"`
	TokenLogoUrl    string `json:"tokenLogoUrl"`
	TokenAddress    string `json:"tokenAddress" dc:"交易数量对应的币种代币地址(转账)"`
	SymbolIn        string `json:"symbolIn" dc:"交易输入币种(兑换)"`
	SymbolOut       string `json:"symbolOut" dc:"交易输出币种(兑换)"`
	InTokenLogoUrl  string `json:"inTokenLogoUrl"`
	OutTokenLogoUrl string `json:"outTokenLogoUrl"`
	TokenAddressIn  string `json:"tokenAddressIn" dc:"交易输入币种代币地址(兑换)"`
	TokenAddressOut string `json:"tokenAddressOut" dc:"交易输出币种代币地址(兑换)"`
	InAmount        string `json:"inAmount" dc:"交易输入数量(兑换)"`
	OutAmount       string `json:"outAmount" dc:"交易输出数量(兑换)"`
	okx.TxHash
}

type TotalValueByAddressReq struct {
	g.Meta           `path:"/total-value-by-address" tags:"钱包" method:"get" summary:"钱包地址总估值"`
	Address          string `json:"address" dc:"具体某条链的账户地址" required:"true"`
	Chains           string `json:"chains" dc:"筛选需要查询交易历史的链，多条链以,分隔" required:"true" d:"501"`
	AssetType        string `json:"assetType" dc:"查询余额类型，默认查所有资产总余额。0：查询所有资产总余额，包括 token 和 defi 资产 1：只查 token 总余额 2：只查 defi 总余额" required:"false"`
	ExcludeRiskToken bool   `json:"excludeRiskToken" dc:"可能存在风险空投代币，可选择是否过滤。默认过滤。true:过滤，false:不过滤" required:"false" d:"true"`
}

type TotalValueByAddressRes model.TotalValueRes

type AllTokenBalancesByAddressReq struct {
	g.Meta  `path:"/all-token-balances-by-address" tags:"钱包" method:"get" summary:"资产明细"`
	Address string `json:"address" dc:"具体某条链的账户地址" required:"true"`
	Chains  string `json:"chains" dc:"筛选需要查询交易历史的链，多条链以,分隔" required:"true" d:"501"`
	Filter  string `json:"filter" dc:"0: 过滤风险空投币 1: 不过滤风险空投币 默认过滤" required:"false"`
}

type AllTokenBalancesByAddressRes struct {
	Data              []model.TokenAsset `json:"data"`
	GasFee            float64            `json:"gasFee" dc:"基础gas费"`
	Earnings          float64            `json:"earnings" dc:"总收益"`
	EarningsRate      int                `json:"earningsRate" dc:"总收益率 最小单位0.01%=1"`
	DailyEarnings     float64            `json:"dailyEarnings" dc:"今日收益"`
	DailyEarningsRate int                `json:"dailyEarningsRate" dc:"今日收益率 最小单位0.01%=1"`
}

type TokenBalancesByAddressReq struct {
	g.Meta         `path:"/token-balances-by-address" tags:"钱包" method:"get" summary:"特定资产明细"`
	Address        string `json:"address" dc:"具体某条链的账户地址" required:"true"`
	Chains         string `json:"chains" dc:"筛选需要查询交易历史的链，多条链以,分隔" required:"true" d:"501"`
	TokenAddresses string `json:"tokenAddresses" dc:"最多支持20个token，用逗号分割" required:"true"`
}

type TokenBalancesByAddressRes []model.TokenAsset
