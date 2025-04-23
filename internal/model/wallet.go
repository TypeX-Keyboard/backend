package model

import "keyboard-api-go/internal/util/chain/okx"

type SupportedChainsRes []okx.Chain

type TransactionsByAddress struct {
	ChainIndex      string   `json:"chainIndex" dc:"chain ID"`
	TxHash          string   `json:"txHash" dc:"transaction hash"`
	TxTime          string   `json:"txTime" dc:"trading hours; Unix timestamps in millisecond format, such as 1597026383085"`
	Action          int      `json:"action" dc:"1 = Receive, 2 = Send, 3 = Buy, 4 = Sell"`
	TxStatus        string   `json:"txStatus" dc:"Transaction Status, Success, Fail, Pending"`
	TokenLogoUrl    string   `json:"tokenLogoUrl"`
	Symbol          string   `json:"symbol" dc:"The currency corresponding to the number of transactions (transfer)"`
	TokenAddress    string   `json:"tokenAddress" dc:"Token address corresponding to the number of transactions (transfer)"`
	Amount          string   `json:"amount" dc:"Number of Transactions (Transfers)"`
	SymbolIn        string   `json:"symbolIn" dc:"Transaction input currency (conversion)"`
	SymbolOut       string   `json:"symbolOut" dc:"Transaction Output Currency (Conversion)"`
	InTokenLogoUrl  string   `json:"inTokenLogoUrl"`
	OutTokenLogoUrl string   `json:"outTokenLogoUrl"`
	TokenAddressIn  string   `json:"tokenAddressIn" dc:"Transaction input currency: token address (exchange)"`
	TokenAddressOut string   `json:"tokenAddressOut" dc:"Transaction Output Currency: Token Address (Exchange)"`
	InAmount        string   `json:"inAmount" dc:"Transaction Input Quantity (Exchange)"`
	OutAmount       string   `json:"outAmount" dc:"Number of Transaction Outputs (Swaps)"`
	Transactions    []TxData `json:"transactions"`
}

type TokenData struct {
	Symbol       string  `json:"symbol" dc:"The currency corresponding to the number of transactions (transfer)"`
	TokenAddress string  `json:"tokenAddress" dc:"Token address corresponding to the number of transactions (transfer)"`
	Amount       float64 `json:"amount" dc:"Number of Transactions (Transfers)"`
}

type TxData struct {
	From         []okx.TxAddress `json:"from" dc:"Transaction Input"`
	To           []okx.TxAddress `json:"to" dc:"Transaction output"`
	TokenAddress string          `json:"tokenAddress" dc:"The contract address of the token"`
	Amount       string          `json:"amount" dc:"Number of transactions"`
	Symbol       string          `json:"symbol" dc:"The currency corresponding to the number of transactions"`
	TxStatus     string          `json:"txStatus" dc:"Transaction Status, Success, Fail, Pending"`
	HitBlacklist bool            `json:"hitBlacklist" dc:"false: is not a blacklist true: is a blacklist"`
}

type TransactionsDetailRes []okx.TxHash

type TotalValueRes []okx.Value

type AllTokenBalancesRes []TokenAsset

type TokenAsset struct {
	okx.TokenAsset
	CostPrice        string  `json:"costPrice" dc:"Cost price"`
	Earning          float64 `json:"earning" dc:"Total Revenue"`
	EarningRate      int     `json:"earningRate" dc:"Yield, 1 means 0.01%"`
	DailyEarning     float64 `json:"dailyEarning" dc:"Daily earnings"`
	DailyEarningRate int     `json:"dailyEarningRate" dc:"Yield, 1 means 0.01%"`
}

type TokenObj struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}
