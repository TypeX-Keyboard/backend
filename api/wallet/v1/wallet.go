package v1

import (
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/util/chain/okx"

	"github.com/gogf/gf/v2/frame/g"
)

type GetSupportedChainsReq struct {
	g.Meta `path:"/supported-chains" tags:"wallet" method:"get" summary:"Query supported chains"`
}

type GetSupportedChainsRes model.SupportedChainsRes

type GetTransactionsByAddressReq struct {
	g.Meta       `path:"/transactions-by-address" tags:"wallet" method:"get" summary:"Query transaction history at the address dimension"`
	Address      string `json:"address" dc:"The account address of a specific chain" required:"true"`
	Chains       string `json:"chains" dc:"Filter the chains that need to query the transaction history, and separate multiple chains" required:"true" d:"501"`
	TokenAddress string `json:"tokenAddress" dc:"Token address 1: Empty represents the main chain token of the corresponding chain. 2: Pass the specific token contract address, which represents the query of the corresponding token. 3: Do not pass, which means querying the main chain coins and all tokens." required:"false"`
	Begin        string `json:"begin" dc:"Start time, query the trading history later than that time. Unix timestamps, expressed in milliseconds" required:"false"`
	End          string `json:"end" dc:"End time, query the transaction history before that time. If both begin and end are not passed, query the transaction history before the current time. Unix timestamps, expressed in milliseconds" required:"false"`
	Cursor       string `json:"cursor" dc:"cursor" required:"false"`
	Limit        string `json:"limit" dc:"The number of returned records is the most recent 20 by default. A maximum of 20 multi-chain queries and 100 single-chain queries are supported" required:"false"`
}

type GetTransactionsByAddressRes struct {
	Cursor          string                        `json:"cursor"`
	TransactionList []model.TransactionsByAddress `json:"transactionList"`
}

type GetTransactionByHashReq struct {
	g.Meta  `path:"/transaction-by-hash" tags:"wallet" method:"get" summary:"Transaction hash to query transaction details"`
	Chain   string `json:"chain" dc:"Filter the chain that needs to query the transaction history" required:"true" d:"501"`
	TxHash  string `json:"txHash" dc:"Transaction hash"`
	Address string `json:"address" dc:"The account address of a specific chain" required:"true"`
}

type GetTransactionByHashRes struct {
	PlatformFee     string `json:"platformFee" dc:"Platform Fees"` // 平台手续费
	Action          int    `json:"action" dc:"1 = Receive, 2 = Send, 3 = Buy, 4 = Sell"`
	TokenLogoUrl    string `json:"tokenLogoUrl"`
	TokenAddress    string `json:"tokenAddress" dc:"Token address corresponding to the number of transactions (transfer)"`
	SymbolIn        string `json:"symbolIn" dc:"Transaction input currency (conversion)"`
	SymbolOut       string `json:"symbolOut" dc:"Transaction Output Currency (Conversion)"`
	InTokenLogoUrl  string `json:"inTokenLogoUrl"`
	OutTokenLogoUrl string `json:"outTokenLogoUrl"`
	TokenAddressIn  string `json:"tokenAddressIn" dc:"Transaction input currency: token address (exchange)"`
	TokenAddressOut string `json:"tokenAddressOut" dc:"Transaction Output Currency: Token Address (Exchange)"`
	InAmount        string `json:"inAmount" dc:"Transaction Input Quantity (Exchange)"`
	OutAmount       string `json:"outAmount" dc:"Number of Transaction Outputs (Swaps)"`
	okx.TxHash
}

type TotalValueByAddressReq struct {
	g.Meta           `path:"/total-value-by-address" tags:"钱包" method:"get" summary:"Total wallet address valuation"`
	Address          string `json:"address" dc:"The account address of a specific chain" required:"true"`
	Chains           string `json:"chains" dc:"Filter the chains that need to query the transaction history, and separate multiple chains" required:"true" d:"501"`
	AssetType        string `json:"assetType" dc:"Query the balance type, and check the total balance of all assets by default. 0: Query the total balance of all assets, including token and DeFi assets 1: Check only the total balance of Token 2: Check only the total balance of DeFi" required:"false"`
	ExcludeRiskToken bool   `json:"excludeRiskToken" dc:"There may be a risk of airdropping tokens, you can choose whether to filter or not. Default filtering. true: filter, false: do not filter" required:"false" d:"true"`
}

type TotalValueByAddressRes model.TotalValueRes

type AllTokenBalancesByAddressReq struct {
	g.Meta  `path:"/all-token-balances-by-address" tags:"wallet" method:"get" summary:"Asset details"`
	Address string `json:"address" dc:"The account address of a specific chain" required:"true"`
	Chains  string `json:"chains" dc:"Filter the chains that need to query the transaction history, and separate multiple chains" required:"true" d:"501"`
	Filter  string `json:"filter" dc:"0: Filter risk airdrop coins 1: Do not filter risk airdrop coins Default filtering" required:"false"`
}

type AllTokenBalancesByAddressRes struct {
	Data              []model.TokenAsset `json:"data"`
	GasFee            float64            `json:"gasFee" dc:"Base gas fee"`
	Earnings          float64            `json:"earnings" dc:"Total Revenue"`
	EarningsRate      int                `json:"earningsRate" dc:"Total Yield: Minimum unit 0.01%=1"`
	DailyEarnings     float64            `json:"dailyEarnings" dc:"Today's earnings"`
	DailyEarningsRate int                `json:"dailyEarningsRate" dc:"Today's Yield Smallest unit 0.01% = 1"`
}

type TokenBalancesByAddressReq struct {
	g.Meta         `path:"/token-balances-by-address" tags:"wallet" method:"get" summary:"Specific asset details"`
	Address        string `json:"address" dc:"The account address of a specific chain" required:"true"`
	Chains         string `json:"chains" dc:"Filter the chains that need to query the transaction history, and separate multiple chains" required:"true" d:"501"`
	TokenAddresses string `json:"tokenAddresses" dc:"A maximum of 20 tokens are supported, separated by commas" required:"true"`
}

type TokenBalancesByAddressRes []model.TokenAsset
