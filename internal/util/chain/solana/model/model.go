package model

import "time"

type Token struct {
	Name     string
	Symbol   string
	Address  string
	Decimals int
}

type From struct {
	Address    string
	PrivateKey string
}

type TransactionRes struct {
	Hex         string
	FromAddress string
	ToAddress   string
	Coin        Token
	Amount      int64
	GasPrice    int64
	GasLimit    int64
}

type TransferObj struct {
	PubKey string
	Amount uint64
}

type RaydiumBaseRes struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Version string `json:"version"`
	Msg     string `json:"msg"`
}

type RaydiumPriorityFeeBaseRes struct {
	RaydiumBaseRes
	Data struct {
		Default PriorityFeeRes `json:"default"`
	} `json:"data"`
}

type PriorityFeeRes struct {
	Vh uint64 `json:"vh"`
	H  uint64 `json:"h"`
	M  uint64 `json:"m"`
}

type QuoteBaseReq struct {
	InputMint  string `json:"inputMint" dc:"Enter the token minting address" binding:"required"`
	OutputMint string `json:"outputMint" dc:"Output token minting address" binding:"required"`
	Amount     uint64 `json:"amount" dc:"The amount to be exchanged, using the smallest unit lamports, 100000000=0.1SOL, inputAmount or outpoutAmount depends on the exchange mode" binding:"required"`
}

type SwapQuoteReq struct {
	QuoteBaseReq
	SlippageBps uint64 `json:"slippageBps" dc:"For example, 0.01% is 1 and 0.03 is 3" binding:"required"`
	PriorityFee uint64 `json:"priorityFee" dc:"The priority transaction fee uses the smallest unit, lamports" binding:"required"`
	UsdValue    uint64 `json:"usdValue" dc:"How much USD is the trade worth" binding:"required"`
	FeeAccount  string `json:"feeAccount" dc:"Expense Account" binding:"required"`
	Mev         bool   `json:"mev"`
}

type RaydiumSwapQuoteEarlyReq struct {
	SwapQuoteReq
	TxVersion string `json:"txVersion" dc:"Transaction version, e.g. LEGACY, V0" binding:"required"`
}
type RaydiumSwapQuoteEarlyRes struct {
	RaydiumBaseRes
	Data struct {
		SwapType             string  `json:"swapType,omitempty"`
		InputMint            string  `json:"inputMint,omitempty"`
		InputAmount          string  `json:"inputAmount,omitempty"`
		OutputMint           string  `json:"outputMint,omitempty"`
		OutputAmount         string  `json:"outputAmount,omitempty"`
		OtherAmountThreshold string  `json:"otherAmountThreshold,omitempty"`
		SlippageBps          uint64  `json:"slippageBps,omitempty"`
		PriceImpactPct       float64 `json:"priceImpactPct,omitempty"`
		ReferrerAmount       string  `json:"referrerAmount,omitempty"`
		RoutePlan            []struct {
			PoolID            string        `json:"poolId,omitempty"`
			InputMint         string        `json:"inputMint,omitempty"`
			OutputMint        string        `json:"outputMint,omitempty"`
			FeeMint           string        `json:"feeMint,omitempty"`
			FeeRate           int           `json:"feeRate,omitempty"`
			FeeAmount         string        `json:"feeAmount,omitempty"`
			RemainingAccounts []interface{} `json:"remainingAccounts"`
			LastPoolPriceX64  string        `json:"lastPoolPriceX64,omitempty"`
		} `json:"routePlan,omitempty"`
	} `json:"data,omitempty"`
}

type RaydiumSwapQuoteLastReq struct {
	ComputeUnitPriceMicroLamports string `json:"computeUnitPriceMicroLamports" dc:"You can set the automatic amount manually or by using the Raydium Priority Fee API via "String(data.data.default.h)". The "h" here stands for high priority. "L" for low, "M" for medium, "H" for high, and "vh" for very high" binding:"required"`
	//SwapResponse                  RaydiumSwapQuoteEarlyRes `json:"swapResponse" dc:"The body of the response returned by the previous request" binding:"required"`
	SwapResponseRaw     interface{} `json:"swapResponse" dc:"The body of the response returned by the previous request" binding:"required"`
	TxVersion           string      `json:"txVersion" dc:"Use 'V0' for versioned transactions and 'LEGACY' for legacy transactions" binding:"required"`
	Wallet              string      `json:"wallet" dc:"pubkey" binding:"required"`
	WrapSol             bool        `json:"wrapSol,omitempty" dc:"It needs to be true to accept SOL as an inputToken"`
	UnwrapSol           bool        `json:"unwrapSol,omitempty" dc:"It needs to be set to true in order to unlock the wSol received as an outputToken"`
	InputTokenAccount   string      `json:"inputTokenAccount,omitempty" dc:"default to ATA "`
	OutpoutTokenAccount string      `json:"outpoutTokenAccount,omitempty" dc:"default to ATA "`
}

type RaydiumSwapQuoteLastRes struct {
	RaydiumBaseRes
	Data []struct {
		Transaction string `json:"transaction"`
	} `json:"data"`
}

type RaydiumMintPriceRes struct {
	RaydiumBaseRes
	Data map[string]*string `json:"data"`
}

type RaydiumMintInfoRes struct {
	RaydiumBaseRes
	Data []RaydiumMintInfoData `json:"data"`
}

type RaydiumMintInfoData struct {
	ChainID    int           `json:"chainId"`
	Address    string        `json:"address"`
	ProgramID  string        `json:"programId"`
	LogoURI    string        `json:"logoURI"`
	Symbol     string        `json:"symbol"`
	Name       string        `json:"name"`
	Decimals   int           `json:"decimals"`
	Tags       []interface{} `json:"tags"`
	Extensions struct {
	} `json:"extensions"`
}

type JupSwapQuoteRes struct {
	InputMint            string      `json:"inputMint"`
	InAmount             string      `json:"inAmount"`
	OutputMint           string      `json:"outputMint"`
	OutAmount            string      `json:"outAmount"`
	OtherAmountThreshold string      `json:"otherAmountThreshold"`
	SwapMode             string      `json:"swapMode"`
	SlippageBps          uint64      `json:"slippageBps"`
	ComputedAutoSlippage uint64      `json:"computedAutoSlippage"`
	PlatformFee          interface{} `json:"platformFee"`
	PriceImpactPct       string      `json:"priceImpactPct"`
	RoutePlan            []struct {
		SwapInfo struct {
			AmmKey     string `json:"ammKey"`
			Label      string `json:"label"`
			InputMint  string `json:"inputMint"`
			OutputMint string `json:"outputMint"`
			InAmount   string `json:"inAmount"`
			OutAmount  string `json:"outAmount"`
			FeeAmount  string `json:"feeAmount"`
			FeeMint    string `json:"feeMint"`
		} `json:"swapInfo"`
		Percent int `json:"percent"`
	} `json:"routePlan"`
	ScoreReport interface{} `json:"scoreReport"`
	ContextSlot int         `json:"contextSlot"`
	TimeTaken   float64     `json:"timeTaken"`
}

type GetSwapRawReq struct {
	QuoteResponse             JupSwapQuoteRes        `json:"quoteResponse"`
	FeeAccount                string                 `json:"feeAccount"`
	UserPublicKey             string                 `json:"userPublicKey"`
	UseSharedAccounts         bool                   `json:"useSharedAccounts "`
	DynamicComputeUnitLimit   bool                   `json:"dynamicComputeUnitLimit"`
	PrioritizationFeeLamports map[string]interface{} `json:"prioritizationFeeLamports"`
	WrapAndUnwrapSol          bool                   `json:"wrapAndUnwrapSol"`
	DynamicSlippage           bool                   `json:"dynamicSlippage"`
}

type GetSwapRawRes struct {
	SwapTransaction           string `json:"swapTransaction"`
	LastValidBlockHeight      int    `json:"lastValidBlockHeight"`
	PrioritizationFeeLamports int    `json:"prioritizationFeeLamports"`
	ComputeUnitLimit          int    `json:"computeUnitLimit"`
	PrioritizationType        struct {
		ComputeBudget struct {
			MicroLamports          int `json:"microLamports"`
			EstimatedMicroLamports int `json:"estimatedMicroLamports"`
		} `json:"computeBudget"`
	} `json:"prioritizationType"`
	SimulationSlot        int         `json:"simulationSlot"`
	DynamicSlippageReport interface{} `json:"dynamicSlippageReport"`
	SimulationError       interface{} `json:"simulationError"`
}

type GetSwapRouterRes struct {
	Wallet      string      `json:"wallet"`
	InputMint   string      `json:"inputMint"`
	InAmount    string      `json:"inAmount"`
	OutputMint  string      `json:"outputMint"`
	OutAmount   string      `json:"outAmount"`
	RoutePlan   []RoutePlan `json:"routePlan,omitempty"`
	TxRaw       string      `json:"txRaw"`
	SlippageBps uint64      `json:"slippageBps"`
	PriorityFee uint64      `json:"priorityFee"`
}

type RoutePlan struct {
	InputMint  string `json:"inputMint,omitempty"`
	OutputMint string `json:"outputMint,omitempty"`
	FeeMint    string `json:"feeMint,omitempty"`
	FeeAmount  string `json:"feeAmount,omitempty"`
}

type JupiterPriceV2Res struct {
	Data      map[string]JupiterPriceV2Data `json:"data"`
	TimeTaken float64                       `json:"timeTaken"`
}

type JupiterPriceV2Data struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Price     string `json:"price"`
	ExtraInfo struct {
		LastSwappedPrice struct {
			LastJupiterSellAt    int    `json:"lastJupiterSellAt"`
			LastJupiterSellPrice string `json:"lastJupiterSellPrice"`
			LastJupiterBuyAt     int    `json:"lastJupiterBuyAt"`
			LastJupiterBuyPrice  string `json:"lastJupiterBuyPrice"`
		} `json:"lastSwappedPrice"`
		QuotedPrice struct {
			BuyPrice  string `json:"buyPrice"`
			BuyAt     int    `json:"buyAt"`
			SellPrice string `json:"sellPrice"`
			SellAt    int    `json:"sellAt"`
		} `json:"quotedPrice"`
		ConfidenceLevel string `json:"confidenceLevel"`
		Depth           struct {
			BuyPriceImpactRatio struct {
				Depth struct {
					Num10   float64 `json:"10"`
					Num100  float64 `json:"100"`
					Num1000 float64 `json:"1000"`
				} `json:"depth"`
				Timestamp int `json:"timestamp"`
			} `json:"buyPriceImpactRatio"`
			SellPriceImpactRatio struct {
				Depth struct {
					Num10   float64 `json:"10"`
					Num100  float64 `json:"100"`
					Num1000 float64 `json:"1000"`
				} `json:"depth"`
				Timestamp int `json:"timestamp"`
			} `json:"sellPriceImpactRatio"`
		} `json:"depth"`
	} `json:"extraInfo"`
}

type TokenInfoRes struct {
	RaydiumMintInfoData
	Price float64 `json:"price"`
}

type GetPoolInfoRes struct {
	Pools []Pool `json:"pools"`
	Total int    `json:"total"`
}

type Pool struct {
	ID        string `json:"id"`
	Chain     string `json:"chain"`
	Dex       string `json:"dex"`
	Type      string `json:"type"`
	BaseAsset struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Symbol       string  `json:"symbol"`
		Icon         string  `json:"icon"`
		Decimals     int     `json:"decimals"`
		CircSupply   float64 `json:"circSupply"`
		TotalSupply  float64 `json:"totalSupply"`
		TokenProgram string  `json:"tokenProgram"`
		FirstPool    struct {
			ID        string    `json:"Id"`
			Dex       string    `json:"Dex"`
			CreatedAt time.Time `json:"CreatedAt"`
		} `json:"firstPool"`
		Fdv      float64 `json:"fdv"`
		Mcap     float64 `json:"mcap"`
		UsdPrice float64 `json:"usdPrice"`
		Stats5M  struct {
			PriceChange       float64 `json:"priceChange"`
			BuyVolume         float64 `json:"buyVolume"`
			SellVolume        float64 `json:"sellVolume"`
			BuyOrganicVolume  float64 `json:"buyOrganicVolume"`
			SellOrganicVolume float64 `json:"sellOrganicVolume"`
			NumBuys           int     `json:"numBuys"`
			NumSells          int     `json:"numSells"`
			NumTraders        int     `json:"numTraders"`
			NumBuyers         int     `json:"numBuyers"`
			NumSellers        int     `json:"numSellers"`
			NumOrganicBuyers  int     `json:"numOrganicBuyers"`
		} `json:"stats5m"`
		Stats1H struct {
			PriceChange       float64 `json:"priceChange"`
			BuyVolume         float64 `json:"buyVolume"`
			SellVolume        float64 `json:"sellVolume"`
			BuyOrganicVolume  float64 `json:"buyOrganicVolume"`
			SellOrganicVolume float64 `json:"sellOrganicVolume"`
			NumBuys           int     `json:"numBuys"`
			NumSells          int     `json:"numSells"`
			NumTraders        int     `json:"numTraders"`
			NumBuyers         int     `json:"numBuyers"`
			NumSellers        int     `json:"numSellers"`
			NumOrganicBuyers  int     `json:"numOrganicBuyers"`
		} `json:"stats1h"`
		Stats6H struct {
			PriceChange       float64 `json:"priceChange"`
			BuyVolume         float64 `json:"buyVolume"`
			SellVolume        float64 `json:"sellVolume"`
			BuyOrganicVolume  float64 `json:"buyOrganicVolume"`
			SellOrganicVolume float64 `json:"sellOrganicVolume"`
			NumBuys           int     `json:"numBuys"`
			NumSells          int     `json:"numSells"`
			NumTraders        int     `json:"numTraders"`
			NumBuyers         int     `json:"numBuyers"`
			NumSellers        int     `json:"numSellers"`
			NumOrganicBuyers  int     `json:"numOrganicBuyers"`
		} `json:"stats6h"`
		Stats24H struct {
			PriceChange       float64 `json:"priceChange"`
			BuyVolume         float64 `json:"buyVolume"`
			SellVolume        float64 `json:"sellVolume"`
			BuyOrganicVolume  float64 `json:"buyOrganicVolume"`
			SellOrganicVolume float64 `json:"sellOrganicVolume"`
			NumBuys           int     `json:"numBuys"`
			NumSells          int     `json:"numSells"`
			NumTraders        int     `json:"numTraders"`
			NumBuyers         int     `json:"numBuyers"`
			NumSellers        int     `json:"numSellers"`
			NumOrganicBuyers  int     `json:"numOrganicBuyers"`
		} `json:"stats24h"`
		Audit struct {
			MintAuthorityDisabled   bool    `json:"mintAuthorityDisabled"`
			FreezeAuthorityDisabled bool    `json:"freezeAuthorityDisabled"`
			TopHoldersPercentage    float64 `json:"topHoldersPercentage"`
		} `json:"audit"`
		OrganicScore      float64  `json:"organicScore"`
		OrganicBuyers24H  int      `json:"organicBuyers24h"`
		OrganicScoreLabel string   `json:"organicScoreLabel"`
		IsVerified        bool     `json:"isVerified"`
		Cexes             []string `json:"cexes"`
	} `json:"baseAsset"`
	CreatedAt time.Time `json:"createdAt"`
	Liquidity float64   `json:"liquidity"`
	Volume24H float64   `json:"volume24h"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type JupTokenInfo struct {
	Address           string      `json:"address"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Decimals          int         `json:"decimals"`
	LogoURI           string      `json:"logoURI"`
	Tags              []string    `json:"tags"`
	DailyVolume       float64     `json:"daily_volume"`
	CreatedAt         time.Time   `json:"created_at"`
	FreezeAuthority   interface{} `json:"freeze_authority"`
	MintAuthority     interface{} `json:"mint_authority"`
	PermanentDelegate interface{} `json:"permanent_delegate"`
	MintedAt          time.Time   `json:"minted_at"`
	Extensions        struct {
		CoingeckoID string `json:"coingeckoId"`
	} `json:"extensions"`
}
