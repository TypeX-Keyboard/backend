package market

import "time"

type MarketRes struct {
	Status Status          `json:"status"`
	Data   map[string]Slug `json:"data"`
}

type Status struct {
	Timestamp    time.Time   `json:"timestamp"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
	Elapsed      int         `json:"elapsed"`
	CreditCount  int         `json:"credit_count"`
	Notice       interface{} `json:"notice"`
}

type Slug struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Symbol         string    `json:"symbol"`
	Slug           string    `json:"slug"`
	NumMarketPairs int       `json:"num_market_pairs"`
	DateAdded      time.Time `json:"date_added"`
	Tags           []struct {
		Slug     string `json:"slug"`
		Name     string `json:"name"`
		Category string `json:"category"`
	} `json:"tags"`
	MaxSupply                     float64          `json:"max_supply"`
	CirculatingSupply             float64          `json:"circulating_supply"`
	TotalSupply                   float64          `json:"total_supply"`
	IsActive                      int              `json:"is_active"`
	InfiniteSupply                bool             `json:"infinite_supply"`
	Platform                      interface{}      `json:"platform"`
	CmcRank                       int              `json:"cmc_rank"`
	IsFiat                        int              `json:"is_fiat"`
	SelfReportedCirculatingSupply interface{}      `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         interface{}      `json:"self_reported_market_cap"`
	TvlRatio                      interface{}      `json:"tvl_ratio"`
	LastUpdated                   time.Time        `json:"last_updated"`
	Quote                         map[string]Quote `json:"quote"`
}

type Quote struct {
	Price                 float64     `json:"price"`
	Volume24H             float64     `json:"volume_24h"`
	VolumeChange24H       float64     `json:"volume_change_24h"`
	PercentChange1H       float64     `json:"percent_change_1h"`
	PercentChange24H      float64     `json:"percent_change_24h"`
	PercentChange7D       float64     `json:"percent_change_7d"`
	PercentChange30D      float64     `json:"percent_change_30d"`
	PercentChange60D      float64     `json:"percent_change_60d"`
	PercentChange90D      float64     `json:"percent_change_90d"`
	MarketCap             float64     `json:"market_cap"`
	MarketCapDominance    float64     `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64     `json:"fully_diluted_market_cap"`
	Tvl                   interface{} `json:"tvl"`
	LastUpdated           time.Time   `json:"last_updated"`
}

type MapRes struct {
	Status Status     `json:"status"`
	Data   []TokenMap `json:"data"`
}

type TokenMap struct {
	ID                  int       `json:"id"`
	Rank                int       `json:"rank"`
	Name                string    `json:"name"`
	Symbol              string    `json:"symbol"`
	Slug                string    `json:"slug"`
	IsActive            int       `json:"is_active"`
	Status              int       `json:"status"`
	FirstHistoricalData time.Time `json:"first_historical_data"`
	LastHistoricalData  time.Time `json:"last_historical_data"`
	Platform            struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Symbol       string `json:"symbol"`
		Slug         string `json:"slug"`
		TokenAddress string `json:"token_address"`
	} `json:"platform"`
}
