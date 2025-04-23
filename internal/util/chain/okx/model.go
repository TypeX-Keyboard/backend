package okx

type Address struct {
	ChainIndex string `json:"chainIndex"`
	Address    string `json:"address"`
}

type CreateWalletReq struct {
	Addresses []Address `json:"addresses"`
}

type CreateWalletRes struct {
	Code string `json:"code"`
	Data []struct {
		Cursor   string `json:"cursor"`
		Accounts []struct {
			AccountID   string `json:"accountId"`
			AccountType string `json:"accountType"`
		} `json:"accounts"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type GetWalletAccountRes struct {
	Code string `json:"code"`
	Data []struct {
		Cursor   string `json:"cursor"`
		Accounts []struct {
			AccountID   string `json:"accountId"`
			AccountType string `json:"accountType"`
		} `json:"accounts"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type PriceRes struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data []PriceData `json:"data"`
}

type PriceData struct {
	ChainIndex   string `json:"chainIndex"`
	TokenAddress string `json:"tokenAddress"`
	Time         string `json:"time"`
	Price        string `json:"price"`
}

type GetPriceInput struct {
	ChainIndex int    `json:"chainIndex"`
	Address    string `json:"tokenAddress"`
}

type Chain struct {
	Name       string `json:"name"`
	LogoUrl    string `json:"logoUrl"`
	ShortName  string `json:"shortName"`
	ChainIndex string `json:"chainIndex"`
}
type ChainRes struct {
	Code string  `json:"code"`
	Msg  string  `json:"msg"`
	Data []Chain `json:"data"`
}

type TxHistoryRes struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []TxData `json:"data"`
}

type TxData struct {
	Cursor          string        `json:"cursor"`
	TransactionList []Transaction `json:"transactionList"`
}

type Transaction struct {
	ChainIndex   string      `json:"chainIndex" dc:"链 ID"`
	TxHash       string      `json:"txHash" dc:"交易 hash"`
	MethodID     string      `json:"methodId" dc:"合约调用函数"`
	Nonce        string      `json:"nonce" dc:"发起者地址发起的第几笔交易"`
	TxTime       string      `json:"txTime" dc:"交易时间；Unix时间戳的毫秒数格式，如 1597026383085"`
	From         []TxAddress `json:"from" dc:"交易输入"`
	To           []TxAddress `json:"to" dc:"交易输出"`
	TokenAddress string      `json:"tokenAddress" dc:"代币的合约地址"`
	Amount       string      `json:"amount" dc:"交易数量"`
	Symbol       string      `json:"symbol" dc:"交易数量对应的币种"`
	TxFee        string      `json:"txFee" dc:"手续费"`
	TxStatus     string      `json:"txStatus" dc:"交易状态、success 成功、fail 失败、pending 等待确认"`
	HitBlacklist bool        `json:"hitBlacklist" dc:"	false：不是黑名单 true：是黑名单"`
	Tag          string      `json:"tag" dc:"黑地址标签类型，包括貔貅盘、网络钓鱼以及合约漏洞等类型。已废弃"`
	IType        string      `json:"itype" dc:"游标"`
}

type TxAddress struct {
	Address string `json:"address" dc:"地址，多签交易时，逗号分隔"`
	Amount  string `json:"amount" dc:"数量"`
}

type TotalValueByAddressRes struct {
	Code string  `json:"code"`
	Msg  string  `json:"msg"`
	Data []Value `json:"data"`
}

type Value struct {
	TotalValue string `json:"totalValue"`
}

type AllTokenBalancesByAddressRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		TokenAssets []TokenAsset `json:"tokenAssets"`
	} `json:"data"`
}

type TokenAsset struct {
	ChainIndex      string `json:"chainIndex" dc:"链唯一标识"`
	TokenAddress    string `json:"tokenAddress" dc:"合约地址"`
	Symbol          string `json:"symbol" dc:"代币简称"`
	Balance         string `json:"balance" dc:"代币数量"`
	TokenPrice      string `json:"tokenPrice" dc:"币种单位价值，以美元计价"`
	TokenType       string `json:"tokenType" dc:"币种类型:1：token 2：铭文"`
	IsRiskToken     bool   `json:"isRiskToken" dc:""`
	TransferAmount  string `json:"transferAmount" dc:"BRC-20、FBRC-20 等铭文资产可直接转账、交易的余额数量，一般称为可转余额。"`
	AvailableAmount string `json:"availableAmount" dc:"BRC-20、FBRC-20 等铭文资产需要完成铭刻操作才可以交易、转账的数量，一般称为可用余额或者待铭刻余额。"`
	RawBalance      string `json:"rawBalance" dc:"true：命中风险空投币 false：未命中风险空投币"`
	Address         string `json:"address" dc:"地址"`
}

type TokenDetailRes struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []TokenDetail `json:"data"`
}

type TokenDetail struct {
	LogoUrl           string              `json:"logoUrl"`
	OfficialWebsite   string              `json:"officialWebsite"`
	SocialUrls        map[string][]string `json:"socialUrls"`
	Decimals          string              `json:"decimals"`
	TokenAddress      string              `json:"tokenAddress"`
	ChainIndex        string              `json:"chainIndex"`
	ChainName         string              `json:"chainName"`
	Symbol            string              `json:"symbol"`
	CirculatingSupply string              `json:"circulatingSupply"`
	MaxSupply         string              `json:"maxSupply"`
	TotalSupply       string              `json:"totalSupply"`
	Volume24h         string              `json:"volume24h"`
	MarketCap         string              `json:"marketCap"`
}

type TokenListRes struct {
	Code string  `json:"code"`
	Msg  string  `json:"msg"`
	Data []Token `json:"data"`
}

type Token struct {
	Decimals             string `json:"decimals"`
	TokenContractAddress string `json:"tokenContractAddress"`
	TokenLogoUrl         string `json:"tokenLogoUrl"`
	TokenName            string `json:"tokenName"`
	TokenSymbol          string `json:"tokenSymbol"`
}

type TxHashRes struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []TxHash `json:"data"`
}

type TxHash struct {
	ChainIndex   string `json:"chainIndex" dc:"链的唯一标识"`
	Height       string `json:"height" dc:"交易发生的区块高度"`
	TxTime       string `json:"txTime" dc:"交易时间；Unix时间戳的毫秒数格式"`
	Txhash       string `json:"txhash" dc:"交易哈希"`
	GasLimit     string `json:"gasLimit" dc:"gas限额"`
	GasUsed      string `json:"gasUsed" dc:"gas消耗"`
	GasPrice     string `json:"gasPrice" dc:"gas价格"`
	TxFee        string `json:"txFee" dc:"交易手续费"`
	Nonce        string `json:"nonce" dc:"nonce"`
	Symbol       string `json:"symbol" dc:"交易数量对应的币种简称"`
	Amount       string `json:"amount" dc:"	交易数量"`
	TxStatus     string `json:"txStatus" dc:"交易状态1:pending 确认中2:success：成功3:fail：失败"`
	MethodId     string `json:"methodId" dc:"合约调用函数"`
	L1OriginHash string `json:"l1OriginHash" dc:"L1执行的交易哈希"`
	FromDetails  []struct {
		Address      string `json:"address" dc:"发送/输入地址"`
		VinIndex     string `json:"vinIndex" dc:"位于当前交易输入的序号"`
		PreVoutIndex string `json:"preVoutIndex" dc:"位于上一笔输出里的序号"`
		TxHash       string `json:"txHash" dc:"交易哈希，和 prevoutIndex 一起唯一确认输入的 UTXO"`
		IsContract   bool   `json:"isContract" dc:"发送地址是否是合约地址 true:是 ；false：否"`
		Amount       string `json:"amount" dc:"交易数量"`
	} `json:"fromDetails" dc:"交易输入详情"`
	ToDetails []struct {
		Address    string `json:"address" dc:"接收/输出地址"`
		VoutIndex  string `json:"voutIndex" dc:"输出的序号"`
		IsContract bool   `json:"isContract" dc:"接收地址是否是合约地址 true:是 ；false：否"`
		Amount     string `json:"amount" dc:"交易数量"`
	} `json:"toDetails" dc:"交易输出详情"`
	InternalTransactionDetails []struct {
		From           string `json:"from" dc:"交易发送方的地址"`
		To             string `json:"to" dc:"交易接受方的地址"`
		IsFromContract bool   `json:"isFromContract" dc:"from地址是否是合约地址"`
		IsToContract   bool   `json:"isToContract" dc:"to地址是否是合约地址"`
		Amount         string `json:"amount" dc:"交易数量"`
		State          string `json:"state" dc:"交易状态"`
	} `json:"internalTransactionDetails" dc:"内部交易详情"`
	TokenTransferDetails []TokenTransferDetail `json:"tokenTransferDetails" dc:"代币交易详情"`
}

type TokenTransferDetail struct {
	Amount               string `json:"amount" dc:"交易数量"`
	From                 string `json:"from" dc:"交易发送方的地址"`
	IsFromContract       bool   `json:"isFromContract" dc:"from地址是否是合约地址"`
	IsToContract         bool   `json:"isToContract" dc:"to地址是否是合约地址"`
	Symbol               string `json:"symbol" dc:"交易代币的简称"`
	To                   string `json:"to" dc:"交易接受方的地址"`
	TokenContractAddress string `json:"tokenContractAddress" dc:"	代币合约地址"`
}
