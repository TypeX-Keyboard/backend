package solana

import (
	"context"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/util/chain/solana/model"
	"time"
)

func (s *solanaClient) SignTx(ctx context.Context, tx *solana.Transaction, privateKey string) error {
	accountFrom, err := s.getPrivateKey(privateKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return accountFrom
			}
			return nil
		},
	)
	if err != nil {
		g.Log().Error(ctx, err)
		return fmt.Errorf("unable to sign transaction: %w", err)
	}
	return nil
}

func (s *solanaClient) TransferSol(ctx context.Context, fromPublicKey string, transferObjs []model.TransferObj) (tx *solana.Transaction, err error) {
	accountFrom, err := s.getPublicKey(fromPublicKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	recent, err := s.rpcClient(ctx).GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	instructions := make([]solana.Instruction, 0)
	for _, transferObj := range transferObjs {
		accountTo, err := s.getPublicKey(transferObj.PubKey)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, err
		}
		instructions = append(instructions, system.NewTransferInstruction(
			transferObj.Amount,
			*accountFrom,
			*accountTo,
		).Build())
	}
	if len(instructions) == 0 {
		err = errors.New("no transfer instruction")
		return
	}
	instructions = s.setPriorityFee(instructions)
	tx, err = solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(*accountFrom),
	)
	return
}

func (s *solanaClient) SendTx(ctx context.Context, tx *solana.Transaction, mev ...bool) (signature solana.Signature, err error) {
	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentConfirmed,
	}
	m := false
	if len(mev) > 0 {
		m = mev[0]
	}
	for range 5 {
		signature, err = s.rpcClient(ctx, m).SendTransactionWithOpts(ctx, tx, opts)
		if err == nil {
			break
		}
		g.Log().Error(ctx, err)
	}
	return
}

func (s *solanaClient) TransferToken(ctx context.Context, fromPublicKey string, toPubKeys []string, amount uint64, tokenAddress string) (tx *solana.Transaction, err error) {
	accountFrom, err := s.getPublicKey(fromPublicKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	tokenMint, err := s.getPublicKey(tokenAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	accountFromToken, _, err := solana.FindAssociatedTokenAddress(*accountFrom, *tokenMint)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	recent, err := s.rpcClient(ctx).GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	instructions := make([]solana.Instruction, 0)
	tokenAccountMap := make(map[string]struct{})
	for _, toPubKey := range toPubKeys {
		accountTo, err := s.getPublicKey(toPubKey)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, err
		}
		tokenAccountTo, _, err := solana.FindAssociatedTokenAddress(*accountTo, *tokenMint)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, err
		}
		if _, ok := tokenAccountMap[tokenAccountTo.String()]; !ok {
			_, err = s.GetAccountInfo(ctx, tokenAccountTo.String())
			if err != nil {
				g.Log().Error(ctx, tokenAccountTo.String(), "tokenAccountTo not exist, create it", err)
				instructions = append(instructions, associatedtokenaccount.NewCreateInstruction(
					*accountFrom,
					*accountTo,
					*tokenMint,
				).Build())
			}
			tokenAccountMap[tokenAccountTo.String()] = struct{}{}
		}
		instructions = append(instructions, token.NewTransferInstruction(
			amount,
			accountFromToken,
			tokenAccountTo,
			*accountFrom,
			nil,
		).Build())
	}
	if len(instructions) == 0 {
		err = errors.New("no token transfer instruction")
		return
	}
	instructions = s.setPriorityFee(instructions)
	tx, err = solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(*accountFrom),
	)
	return
}

func (s *solanaClient) GetTransaction(ctx context.Context, txSig solana.Signature) (out *rpc.GetTransactionResult, err error) {
	var maxSupportedTransactionVersion uint64 = 0
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 3)
		out, err = s.rpcClient(ctx).GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			MaxSupportedTransactionVersion: &maxSupportedTransactionVersion,
		})
		if err == nil {
			return out, nil
		}
	}
	return nil, err
}

func (s *solanaClient) GetSignatureStatuses(ctx context.Context, txSig solana.Signature) (out *rpc.GetSignatureStatusesResult, err error) {
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 3)
		out, err = s.rpcClient(ctx).GetSignatureStatuses(ctx, true, txSig)
		if err == nil && out.Value[0] != nil {
			return out, nil
		}
		g.Log().Error(ctx, err)
	}
	return nil, err
}
