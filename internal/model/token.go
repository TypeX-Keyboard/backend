package model

import (
	"keyboard-api-go/internal/util/chain/solana/model"
)

type TokenDetail model.Pool

type CheckAddress struct {
	AddressType int                `json:"address_type" dc:"1:token，2:wallet，3:unknown"`
	SolAmount   uint64             `json:"sol_amount" dc:"address_type=2 has a value; The amount of SOL held"`
	Holds       map[string]Balance `json:"holds" dc:"address_type=2 has a value; List of tokens held, the key is tokenAccount"`
}

type Balance struct {
	TokenAddress string `json:"token_address" dc:"tokenAddress"`
	Amount       uint64 `json:"amount" dc:"amount"`
}
