package solana

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/util/str"
)

type solanaClient struct {
	rpcClient func(ctx context.Context, mev ...bool) *rpc.Client
}

type SolanaConf struct {
	Address    string
	PrivateKey string
}

var rpcPoints = []string{
	rpc.MainNetBeta_RPC,
}

var solanaCli *solanaClient

func InitSolana(ctx context.Context, conf SolanaConf) error {
	// Create a new RPC client:
	var rpcClient = func(ctx context.Context, mev ...bool) *rpc.Client {
		if len(mev) > 0 {
			if mev[0] {
				return rpc.New("https://mainnet.block-engine.jito.wtf/api/v1/transactions")
			}
		}
		randomValue := str.PopElement(&rpcPoints)
		g.Log().Infof(ctx, "solana rpc use point: %v", randomValue)
		return rpc.New(randomValue)
	}

	// Create a new WS client (used for confirming transactions)
	solanaCli = &solanaClient{
		rpcClient: rpcClient,
	}
	from.Address = conf.Address
	from.PrivateKey = conf.PrivateKey
	return nil
}

func GetSolana() *solanaClient {
	return solanaCli
}

func (s *solanaClient) GetRpcURL(ctx context.Context) []string {
	return rpcPoints
}

func (s *solanaClient) setPriorityFee(instructions []solana.Instruction) []solana.Instruction {
	instructions = append([]solana.Instruction{computebudget.NewSetComputeUnitPriceInstruction(300).Build()}, instructions...)
	instructions = append([]solana.Instruction{computebudget.NewSetComputeUnitLimitInstruction(50_000).Build()}, instructions...)
	return instructions
}

func (s *solanaClient) getPrivateKey(priKey string) (*solana.PrivateKey, error) {
	privateKey, err := solana.PrivateKeyFromBase58(priKey)
	if err != nil {
		return nil, err
	}
	return &privateKey, nil
}

func (s *solanaClient) getPublicKey(pubKey string) (*solana.PublicKey, error) {
	publicKey, err := solana.PublicKeyFromBase58(pubKey)
	if err != nil {
		return nil, err
	}
	return &publicKey, nil
}

func (s *solanaClient) GetSOLBalance(ctx context.Context, pubKey string) (uint64, error) {
	publicKey, err := s.getPublicKey(pubKey)
	if err != nil {
		return 0, err
	}
	out, err := s.rpcClient(ctx).GetBalance(
		ctx,
		*publicKey,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return 0, err
	}
	return out.Value, nil
}

func (s *solanaClient) GetTokenBalance(ctx context.Context, pubKey string) (map[string]token.Account, error) {
	publicKey, err := s.getPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	out, err := s.rpcClient(ctx).GetTokenAccountsByOwner(
		ctx,
		*publicKey,
		&rpc.GetTokenAccountsConfig{ProgramId: &solana.TokenProgramID},
		&rpc.GetTokenAccountsOpts{Encoding: solana.EncodingBase64Zstd},
	)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	if out == nil {
		return nil, nil
	}
	accounts := make(map[string]token.Account)
	for _, account := range out.Value {
		var tokAcc token.Account
		data := account.Account.Data.GetBinary()
		dec := bin.NewBinDecoder(data)
		err := dec.Decode(&tokAcc)
		if err != nil {
			g.Log().Error(ctx, "GetTokenBalance", err)
			continue
		}
		accounts[account.Pubkey.String()] = tokAcc
	}
	return accounts, nil
}

func (s *solanaClient) GetTokenSupply(ctx context.Context, mintAddress string) (out *rpc.GetTokenSupplyResult, err error) {
	mint, err := s.getPublicKey(mintAddress)
	if err != nil {
		return
	}
	return s.rpcClient(ctx).GetTokenSupply(ctx, *mint, rpc.CommitmentConfirmed)
}

func (s *solanaClient) SignedTxBase64(ctx context.Context, input string) (string, error) {
	accountFrom, err := s.getPrivateKey(from.PrivateKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return "", err
	}
	var transaction solana.Transaction
	err = transaction.UnmarshalBase64(input)
	if err != nil {
		return "", err
	}
	_, err = transaction.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return accountFrom
			}
			return nil
		},
	)
	// Encode the signed transaction as base64
	signedTxBase64, err := transaction.ToBase64()
	g.Log().Infof(ctx, "signedTxBase64 %s", signedTxBase64)
	return signedTxBase64, nil
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *solanaClient) ParseRPCError(ctx context.Context, orgErr error) (rpcError RPCError) {
	err := json.Unmarshal([]byte(orgErr.Error()), &rpcError)
	if err != nil {
		g.Log().Error(ctx, "Error parsing JSON:", orgErr)
		g.Log().Error(ctx, "Error parsing JSON:", err)
		return
	}
	return
}

func (s *solanaClient) CreateMint(ctx context.Context) {
	accountFrom, err := s.getPrivateKey(from.PrivateKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	recent, err := s.rpcClient(ctx).GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// Create a fee payer and Mint account
	feePayer := accountFrom
	alice := accountFrom
	mint := solana.NewWallet()
	// Get the minimum rent waiver balance
	minBalance, err := s.rpcClient(ctx).GetMinimumBalanceForRentExemption(ctx, token.MINT_SIZE, rpc.CommitmentFinalized)
	if err != nil {
		g.Log().Errorf(ctx, "failed to get minimum balance for rent exemption: %v", err)
		return
	}

	// Create a transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			// Create a Mint account
			system.NewCreateAccountInstruction(
				minBalance,
				token.MINT_SIZE,
				token.ProgramID,
				mint.PublicKey(),
				feePayer.PublicKey(),
			).Build(),

			// Initialize the Mint account
			token.NewInitializeMintInstruction(
				6, // Number of decimal places
				alice.PublicKey(),
				solana.PublicKey{},
				mint.PublicKey(),
				token.ProgramID,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(feePayer.PublicKey()),
	)
	if err != nil {
		g.Log().Errorf(ctx, "failed to create transaction: %v", err)
		return
	}
	spew.Dump(tx)
	spew.Dump(mint.PublicKey())
	spew.Dump(mint.PrivateKey)
}
