package chain

import (
	"context"
	"keyboard-api-go/internal/util/chain/solana"
)

type ChainConf struct {
	SolConf solana.SolanaConf
}

func InitChain(ctx context.Context, conf ChainConf) error {
	if err := solana.InitSolana(ctx, conf.SolConf); err != nil {
		return err
	}
	return nil
}
