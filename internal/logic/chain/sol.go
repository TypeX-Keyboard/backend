package chain

import (
	"context"
	"fmt"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	utilSolana "keyboard-api-go/internal/util/chain/solana"
	model2 "keyboard-api-go/internal/util/chain/solana/model"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func (s *sChain) SolTransfer(ctx context.Context, from, to, token string, amount uint64) (res *model.TransferRes, err error) {
	var tx *solana.Transaction
	if token == utilSolana.SolTokens[utilSolana.SOL].Address {
		tx, err = utilSolana.GetSolana().TransferSol(ctx, from, []model2.TransferObj{
			{
				PubKey: to,
				Amount: amount,
			},
		})
		if err != nil {
			return nil, err
		}
	} else {
		tx, err = utilSolana.GetSolana().TransferToken(ctx, from, []string{to}, amount, token)
		if err != nil {
			return nil, err
		}
	}
	res = &model.TransferRes{
		Program: tx,
		TxRaw:   tx.MustToBase64(),
	}
	return res, nil
}

func (s *sChain) SolPriorityFee(ctx context.Context) (model.PriorityFeeRes, error) {
	fee, err := utilSolana.NewDex(utilSolana.Jupiter).GetPriorityFee(ctx)
	if err != nil {
		return model.PriorityFeeRes{}, err
	}
	return model.PriorityFeeRes(fee), err
}

func (s *sChain) SolSlippageBps(ctx context.Context, tokenIn, tokenOut string, amount uint64) (uint64, error) {
	return utilSolana.NewDex(utilSolana.Jupiter).GetSlippageBps(ctx, model2.QuoteBaseReq{
		InputMint:  tokenIn,
		OutputMint: tokenOut,
		Amount:     amount,
	})
}

func (s *sChain) SolSwapRoute(ctx context.Context, adaptor, signAddress, tokenIn, tokenOut string, amountIn, slippageBps, priorityFee uint64, mev bool) (*model.SwapRouterRes, error) {
	res := &model.SwapRouterRes{}
	input := model2.SwapQuoteReq{
		QuoteBaseReq: model2.QuoteBaseReq{
			InputMint:  tokenIn,
			OutputMint: tokenOut,
			Amount:     amountIn,
		},
		SlippageBps: slippageBps,
		FeeAccount:  consts.FeeAccountATA,
		PriorityFee: priorityFee,
		Mev:         mev,
	}
	router, err := utilSolana.NewDex(utilSolana.Adaptor(adaptor)).GetSwapRaw(ctx, input, signAddress)
	if err != nil {
		return res, err
	}
	res.GetSwapRouterRes = router
	return res, nil
}

func (s *sChain) SolGetMetaData(ctx context.Context, adaptor, tokens string) (res model.TokenInfoRes, err error) {
	idex := utilSolana.NewDex(utilSolana.Adaptor(adaptor))
	return idex.GetTokenInfo(ctx, tokens)
}

func (s *sChain) SolSendTx(ctx context.Context, txRaw, signAddress string) (string, error) {
	tx, err := solana.TransactionFromBase64(txRaw)
	if err != nil {
		return "", err
	}
	publicKey, err := solana.PublicKeyFromBase58(signAddress)
	if err != nil {
		return "", err
	}
	ok := tx.IsSigner(publicKey)
	if !ok {
		return "", fmt.Errorf("signAddress is not signer")
	}
	signature, err := utilSolana.GetSolana().SendTx(ctx, tx)
	if err != nil {
		return "", err
	}
	return signature.String(), nil
}

func (s *sChain) SolGetTxInfo(ctx context.Context, hex string) (out *rpc.GetTransactionResult, err error) {
	signature, err := solana.SignatureFromBase58(hex)
	if err != nil {
		return nil, err
	}
	return utilSolana.GetSolana().GetTransaction(ctx, signature)
}

func (s *sChain) SolRpcURL(ctx context.Context) []string {
	return utilSolana.GetSolana().GetRpcURL(ctx)
}
