package solana

import (
	"context"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/util/chain/solana/model"
	"log"
)

var from model.From

func (s *solanaClient) GetAccountInfo(ctx context.Context, pubKey string) (*rpc.GetAccountInfoResult, error) {
	key, err := s.getPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return solanaCli.rpcClient(ctx).GetAccountInfo(ctx, *key)
}

func (s *solanaClient) GetFrom(ctx context.Context) (*solana.PrivateKey, error) {
	privateKey, err := solana.PrivateKeyFromBase58(from.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &privateKey, nil
}

func (s *solanaClient) GenTokenAccount(ctx context.Context, payerAddress, walletAddress, mintAddress string) (tx *solana.Transaction, err error) {
	recent, err := s.rpcClient(ctx).GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	payer, err := s.getPublicKey(payerAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	wallet, err := s.getPublicKey(walletAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	mint, err := s.getPublicKey(mintAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	instructions := make([]solana.Instruction, 0)
	createAccountInstruction := associatedtokenaccount.NewCreateInstruction(
		*payer,
		*wallet,
		*mint,
	).Build()
	instructions = append(instructions, createAccountInstruction)
	instructions = s.setPriorityFee(instructions)
	tx, err = solana.NewTransaction(instructions, recent.Value.Blockhash, solana.TransactionPayer(*payer))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}

func (s *solanaClient) DelTokenAccount(ctx context.Context, walletAddress string, mintAddress []string) (tx *solana.Transaction, err error) {
	recent, err := s.rpcClient(ctx).GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	accountFrom, err := s.getPublicKey(walletAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	balanceMap, err := s.GetTokenBalance(ctx, walletAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	instructions := make([]solana.Instruction, 0)
	mintArr := garray.NewStrArrayFrom(mintAddress)
	for account, info := range balanceMap {
		g.Dump(info.Mint.String())
		if mintArr.Contains(info.Mint.String()) {
			tokenAccountAddress := account
			accountInfo := info
			tokenAccount, err := s.getPublicKey(tokenAccountAddress)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}
			if accountInfo.Amount > 0 {
				burnInstruction := token.NewBurnInstruction(accountInfo.Amount, *tokenAccount, accountInfo.Mint, accountInfo.Owner, nil).Build()
				instructions = append(instructions, burnInstruction)
			}
			closeAccountInstruction := token.NewCloseAccountInstruction(*tokenAccount, accountInfo.Owner, accountInfo.Owner, nil).Build()
			instructions = append(instructions, closeAccountInstruction)
		}
	}
	if len(instructions) == 0 {
		err = errors.New("no token account to delete")
		return
	}
	instructions = s.setPriorityFee(instructions)
	tx, err = solana.NewTransaction(instructions, recent.Value.Blockhash, solana.TransactionPayer(*accountFrom))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}

func (s *solanaClient) GetSignaturesForAddress(ctx context.Context, address string) (out []*rpc.TransactionSignature, err error) {
	accountFrom, err := s.getPublicKey(address)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return s.rpcClient(ctx).GetSignaturesForAddress(ctx, *accountFrom)
}

func (s *solanaClient) GetProgramAccounts(ctx context.Context, tokenAddress string) {
	tokenMint, err := solana.PublicKeyFromBase58(tokenAddress)
	if err != nil {
		log.Fatalf("Invalid token mint address: %v", err)
	}
	resp, err := s.rpcClient(ctx).GetProgramAccountsWithOpts(
		ctx,
		solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"),
		&rpc.GetProgramAccountsOpts{
			Encoding: solana.EncodingBase64,
			Filters: []rpc.RPCFilter{
				{
					DataSize: 165,
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 0,
						Bytes:  tokenMint.Bytes(),
					},
				},
			},
		},
	)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 使用 map 统计唯一持有者地址
	holderSet := make(map[string]struct{})

	for _, account := range resp {
		// 提取账户数据
		accountPubKey := account.Pubkey.String()
		holderSet[accountPubKey] = struct{}{}
	}
	spew.Dump(len(holderSet))
	spew.Dump(len(resp)) // NOTE: this can generate a lot of output
}
