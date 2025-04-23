package solana

import "keyboard-api-go/internal/util/chain/solana/model"

const (
	SOL  = "SOL"
	USDT = "USDT"
)

var SolTokens = map[string]model.Token{
	SOL: {
		Name:     "Solana",
		Symbol:   "Sol",
		Decimals: 9,
		Address:  "So11111111111111111111111111111111111111112",
	},
	USDT: {
		Name:     "Tether",
		Symbol:   "USDT",
		Decimals: 6,
		Address:  "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",
	},
}
