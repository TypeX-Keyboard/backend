package token

import (
	"context"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/chain/solana"
)

type sToken struct {
}

func init() {
	service.RegisterToken(New())
}

func New() service.IToken {
	return &sToken{}
}

func (s *sToken) CheckAddress(ctx context.Context, address string) (*model.CheckAddress, error) {
	info, err := solana.GetSolana().GetAccountInfo(ctx, address)
	if err != nil {
		g.Log().Error(ctx, address)
		g.Log().Error(ctx, err)
		_, err := solanago.PublicKeyFromBase58(address)
		if err != nil {
			return &model.CheckAddress{AddressType: 3}, nil
		}
		return &model.CheckAddress{AddressType: 2}, nil
	}
	switch info.Value.Owner.String() {
	case consts.AccountProgram:
		res := model.CheckAddress{AddressType: 2}
		holds := make(map[string]model.Balance)
		sol, err := solana.GetSolana().GetSOLBalance(ctx, address)
		if err != nil {
			return nil, err
		}
		res.SolAmount = sol
		holds[address] = model.Balance{TokenAddress: consts.SolAddress, Amount: sol}
		tokenBalances, err := solana.GetSolana().GetTokenBalance(ctx, address)
		if err != nil {
			return nil, err
		}
		for k, account := range tokenBalances {
			if account.Amount == 0 {
				continue
			}
			holds[k] = model.Balance{TokenAddress: account.Mint.String(), Amount: account.Amount}
		}
		res.Holds = holds
		return &res, nil
	case consts.Token2022Program:
		return &model.CheckAddress{AddressType: 1}, nil
	case consts.TokenProgram:
		return &model.CheckAddress{AddressType: 1}, nil
	default:
		return &model.CheckAddress{AddressType: 3}, nil
	}
}
