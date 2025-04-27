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
	ChainIndex   string      `json:"chainIndex" dc:"chain ID"`
	TxHash       string      `json:"txHash" dc:"transaction hash"`
	MethodID     string      `json:"methodId" dc:"The contract calls the function"`
	Nonce        string      `json:"nonce" dc:"The first few transactions initiated by the originator's address"`
	TxTime       string      `json:"txTime" dc:"trading hours; Unix timestamps in millisecond format, such as 1597026383085"`
	From         []TxAddress `json:"from" dc:"Transaction Input"`
	To           []TxAddress `json:"to" dc:"Transaction output"`
	TokenAddress string      `json:"tokenAddress" dc:"The contract address of the token"`
	Amount       string      `json:"amount" dc:"Number of transactions"`
	Symbol       string      `json:"symbol" dc:"The currency corresponding to the number of transactions"`
	TxFee        string      `json:"txFee" dc:"Premium"`
	TxStatus     string      `json:"txStatus" dc:"Transaction Status, Success, Fail, Pending"`
	HitBlacklist bool        `json:"hitBlacklist" dc:"	false: is not a blacklist true: is a blacklist"`
	Tag          string      `json:"tag" dc:"Types of black address labels, including Pixiu disks, phishing, and contract vulnerabilities. Deprecated"`
	IType        string      `json:"itype" dc:"cursor"`
}

type TxAddress struct {
	Address string `json:"address" dc:"addresses, multisig transactions, comma separated"`
	Amount  string `json:"amount" dc:"number"`
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
	ChainIndex      string `json:"chainIndex" dc:"Chain unique identification"`
	TokenAddress    string `json:"tokenAddress" dc:"Contract address"`
	Symbol          string `json:"symbol" dc:"Token abbreviation"`
	Balance         string `json:"balance" dc:"Number of tokens"`
	TokenPrice      string `json:"tokenPrice" dc:"The value of the currency unit, denominated in USD"`
	TokenType       string `json:"tokenType" dc:"Token type: 1: token 2: inscription"`
	IsRiskToken     bool   `json:"isRiskToken" dc:""`
	TransferAmount  string `json:"transferAmount" dc:"The number of balances of inscription assets such as BRC-20 and FBRC-20 that can be directly transferred and traded is generally referred to as the transferable balance."`
	AvailableAmount string `json:"availableAmount" dc:"The amount of inscription assets such as BRC-20 and FBRC-20 that need to be engraved before they can be traded or transferred, generally referred to as the available balance or the balance to be inscribed."`
	RawBalance      string `json:"rawBalance" dc:"true: the risky airdrop is hit false: the risky airdrop is missed"`
	Address         string `json:"address" dc:"address"`
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
	ChainIndex   string `json:"chainIndex" dc:"The unique identifier of the chain"`
	Height       string `json:"height" dc:"The height of the block where the transaction occurred"`
	TxTime       string `json:"txTime" dc:"trading hours; Unix timestamps in millisecond format"`
	Txhash       string `json:"txhash" dc:"Transaction hash"`
	GasLimit     string `json:"gasLimit" dc:"Gas limits"`
	GasUsed      string `json:"gasUsed" dc:"Gas consumption"`
	GasPrice     string `json:"gasPrice" dc:"Gas price"`
	TxFee        string `json:"txFee" dc:"Transaction fees"`
	Nonce        string `json:"nonce" dc:"nonce"`
	Symbol       string `json:"symbol" dc:"The abbreviation of the currency corresponding to the number of transactions"`
	Amount       string `json:"amount" dc:"	Number of transactions"`
	TxStatus     string `json:"txStatus" dc:"Transaction status 1: pendingConfirmed2: success: succeeded 3: fail: failed"`
	MethodId     string `json:"methodId" dc:"The contract calls the function"`
	L1OriginHash string `json:"l1OriginHash" dc:"The hash of the transaction executed by L1"`
	FromDetails  []struct {
		Address      string `json:"address" dc:"Send/enter address"`
		VinIndex     string `json:"vinIndex" dc:"The serial number that is located in the current transaction input"`
		PreVoutIndex string `json:"preVoutIndex" dc:"The serial number located in the previous output"`
		TxHash       string `json:"txHash" dc:"The transaction hash, together with the prevoutIndex, uniquely confirms the UTXO of the input"`
		IsContract   bool   `json:"isContract" dc:"Whether the sending address is the contract address true: Yes; false: No"`
		Amount       string `json:"amount" dc:"Number of transactions"`
	} `json:"fromDetails" dc:"Transaction input details"`
	ToDetails []struct {
		Address    string `json:"address" dc:"Receive/Output Address"`
		VoutIndex  string `json:"voutIndex" dc:"The ordinal number of the output"`
		IsContract bool   `json:"isContract" dc:"Whether the receiving address is the contract address true: Yes; false: No"`
		Amount     string `json:"amount" dc:"Number of transactions"`
	} `json:"toDetails" dc:"Details of the transaction output"`
	InternalTransactionDetails []struct {
		From           string `json:"from" dc:"The address of the sender of the transaction"`
		To             string `json:"to" dc:"The address of the party to whom the transaction was accepted"`
		IsFromContract bool   `json:"isFromContract" dc:"Whether the from address is the contract address"`
		IsToContract   bool   `json:"isToContract" dc:"Whether the to address is the contract address"`
		Amount         string `json:"amount" dc:"Number of transactions"`
		State          string `json:"state" dc:"Transaction status"`
	} `json:"internalTransactionDetails" dc:"Insider transaction details"`
	TokenTransferDetails []TokenTransferDetail `json:"tokenTransferDetails" dc:"Token trading details"`
}

type TokenTransferDetail struct {
	Amount               string `json:"amount" dc:"Number of transactions"`
	From                 string `json:"from" dc:"The address of the sender of the transaction"`
	IsFromContract       bool   `json:"isFromContract" dc:"Whether the from address is the contract address"`
	IsToContract         bool   `json:"isToContract" dc:"Whether the to address is the contract address"`
	Symbol               string `json:"symbol" dc:"The short name for the trading token"`
	To                   string `json:"to" dc:"The address of the party to whom the transaction was accepted"`
	TokenContractAddress string `json:"tokenContractAddress" dc:"	Token contract address"`
}
