package solana

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/shopspring/decimal"
	"github.com/tyler-smith/go-bip39"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"keyboard-api-go/internal/util/chain/solana/model"
	"log"
	"strings"
	"testing"
	"time"
)

var conf = SolanaConf{
	Address:    "",
	PrivateKey: "",
}

func init() {
	ctx := gctx.GetInitCtx()
	if err := InitSolana(ctx, conf); err != nil {
		panic(err)
	}
}

func TestSolana(t *testing.T) {
	ctx := gctx.New()
	balance, err := GetSolana().GetSOLBalance(ctx, conf.Address)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(balance)
	accounts, err := GetSolana().GetTokenBalance(ctx, "7BipTb8pQh7xbjioxjzdCGDMV1AJspjnjMXq1e5HjvZN")
	if err != nil {
		t.Error(err)
	}
	spew.Dump(accounts)
}

func TestGetTokenSupply(t *testing.T) {
	ctx := gctx.New()
	out, err := GetSolana().GetTokenSupply(ctx, "6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN")
	if err != nil {
		t.Error(err)
	}
	spew.Dump(out)
}

func TestSol(t *testing.T) {
	ctx := gctx.New()
	owner, err := GetSolana().getPublicKey("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	if err != nil {
		t.Error(err)
		return
	}
	token, err := GetSolana().getPublicKey("6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN")
	if err != nil {
		t.Error(err)
		return
	}
	out, err := GetSolana().rpcClient(ctx).GetTokenAccountsByOwner(ctx, *token, &rpc.GetTokenAccountsConfig{
		ProgramId: owner,
	},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingBase64Zstd,
		})
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(out)
}

func TestGetSignaturesForAddress(t *testing.T) {
	ctx := gctx.New()
	out, err := GetSolana().GetSignaturesForAddress(ctx, conf.Address)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(out)
}

func TestGenTokenAccount(t *testing.T) {
	ctx := gctx.New()
	walletAddress := "Eign4g779oTaexjVnQt6gcQh3SXa1t1sRca4ajrNbvFK"
	token := "So11111111111111111111111111111111111111112"
	tx, err := GetSolana().GenTokenAccount(ctx, conf.Address, walletAddress, token)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(tx)
	err = GetSolana().SignTx(ctx, tx, from.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := GetSolana().SendTx(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(signature)
	out, err := GetSolana().GetTransaction(ctx, signature)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(out.Meta.Status)
}

func TestDelTokenAccount(t *testing.T) {
	ctx := gctx.New()
	accounts, err := GetSolana().GetTokenBalance(ctx, conf.Address)
	if err != nil {
		t.Error(err)
	}
	ray := dex.NewRaydium()
	tokens := make([]string, 0)
	for _, account := range accounts {
		tokens = append(tokens, account.Mint.String())
	}
	prices, err := ray.MintPrice(ctx, tokens)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(prices)
	delToken := make([]string, 0)
	for _, account := range accounts {
		delToken = append(delToken, account.Mint.String())
	}
	g.Dump(delToken)
	tx, err := GetSolana().DelTokenAccount(ctx, conf.Address, delToken)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(tx)
	err = GetSolana().SignTx(ctx, tx, from.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := GetSolana().SendTx(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(signature)
	out, err := GetSolana().GetTransaction(ctx, signature)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(out.Meta.Status)
}

func TestSolanaTransfer(t *testing.T) {
	ctx := gctx.New()
	tx, err := GetSolana().TransferToken(ctx, from.Address, []string{
		"GM4X9e6sCXa9JkvPzXWFyWKRzekcAhchtD18wsHSvG6c",
	}, 100000, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	if err != nil {
		t.Error(err)
	}
	err = GetSolana().SignTx(ctx, tx, from.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	signature, err := GetSolana().SendTx(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(signature)
	out, err := GetSolana().GetTransaction(ctx, signature)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(out.Meta.Status)
}

func TestSolanaTransferSol(t *testing.T) {
	ctx := gctx.New()
	tx, err := GetSolana().TransferSol(ctx, from.Address, []model.TransferObj{
		{PubKey: "CNuw5gMrXzN7sTNnvPy5W8TfzcU7HS2f5k75czCfo8Ht", Amount: 1},
		{PubKey: "9AgLkFvwDwyA5ZHeAUjZveshEwgpGK4krKXqh5jR2Wvg", Amount: 1},
	})
	if err != nil {
		t.Error(err)
	}
	err = GetSolana().SignTx(ctx, tx, from.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	signature, err := GetSolana().SendTx(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(signature)
	out, err := GetSolana().GetTransaction(ctx, signature)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(out.Meta.Status)
}

func TestSolanaClient_GenAccount(t *testing.T) {
	// 生成熵
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(entropy)
	// 生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("生成的助记词:", mnemonic)
	seed := bip39.NewSeed(mnemonic, "")
	priv := ed25519.NewKeyFromSeed(seed[:32])
	solPrivateKey := solana.PrivateKey(priv)
	fmt.Println("private key:", solPrivateKey.String())
	fmt.Println("public key:", solPrivateKey.PublicKey().String())
	fmt.Println("public IsOnCurve:", solPrivateKey.PublicKey().IsOnCurve())
}

func TestGetAccountInfo(t *testing.T) {
	ctx := gctx.New()
	//FJ5f2fiawX4983mZ5AptFv72vJrtZBzgUxStjy1J3zCE
	// account
	key := "6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN"
	info, err := GetSolana().GetAccountInfo(ctx, key)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(info)
	publicKey, err := GetSolana().getPublicKey(key)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(publicKey)
	// token
	//key = "6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN"
	//info, err = GetSolana().GetAccountInfo(ctx, key)
	//if err != nil {
	//	t.Error(err)
	//}
	//spew.Dump(info)
	//var mint token.Mint
	//// Account{}.Data.GetBinary() returns the *decoded* binary data
	//// regardless the original encoding (it can handle them all).
	//err = bin.NewBinDecoder(info.GetBinary()).Decode(&mint)
	//if err != nil {
	//	panic(err)
	//}
	//spew.Dump(mint)
	//// program
	//key = "JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"
	//info, err = GetSolana().GetAccountInfo(ctx, key)
	//if err != nil {
	//	t.Error(err)
	//}
	//spew.Dump(info)
}
func TestGetTransaction(t *testing.T) {
	//FJ5f2fiawX4983mZ5AptFv72vJrtZBzgUxStjy1J3zCE
	hash := "2MCKPoPpKUnRyDFa5gLYunSUsiisXrNKqbzt15XfPmd6tpqAWTPDcK2DwE2eaE1gXzG44rmC9xhTFpRYiEjxYtRi"
	ctx := gctx.New()
	txSig := solana.MustSignatureFromBase58(hash)
	out, err := GetSolana().GetTransaction(ctx, txSig)
	if err != nil {
		t.Error(err)
	}
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
	if err != nil {
		panic(err)
	}
	g.Dump(tx)
}

func TestBip44(t *testing.T) {
	mnemonic := "pill tomorrow foster good walnut borrow virtual kick shift mutual shoe scatter"
	fmt.Println("\n生成的助记词:")
	fmt.Println("  助记词:", mnemonic)

	seed := bip39.NewSeed(mnemonic, "")
	fmt.Printf("完整种子 (Hex): %x\n", seed)

	hd, err := FromMasterSeed(seed)
	if err != nil {
		log.Fatal(err)
	}

	// 打印链码和私钥
	fmt.Printf("chainCode (Base58): %s\n", base58.Encode(hd.ChainCode))
	fmt.Printf("privateKey (Base58): %s\n", base58.Encode(hd.PrivateKey))
	address := make([]string, 0)
	// 派生子密钥
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("m/44'/501'/%d'/0'", i)
		derived, err := hd.Derive(path)
		if err != nil {
			continue
		}

		// 注意：这里需要实现 Keypair 结构体和相关方法
		priv := ed25519.NewKeyFromSeed(derived.PrivateKey)
		solPrivateKey := solana.PrivateKey(priv)
		fmt.Printf("%s => 公钥：%s\n", path, solPrivateKey.PublicKey())
		fmt.Printf("%s => 私钥：%s\n", path, solPrivateKey.String())
		fmt.Printf("%s => 公钥验证：%v\n", path, solPrivateKey.PublicKey().IsOnCurve())
		poload := []byte("Hello, Solana!")
		sign, err := solPrivateKey.Sign(poload)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("sign:", sign)
		verify := sign.Verify(solPrivateKey.PublicKey(), poload)
		fmt.Println("签名验证verify:", verify)
		if verify {
			address = append(address, solPrivateKey.PublicKey().String())
		}
	}
	for _, s := range address {
		g.Dump(s)
	}
}

func TestTransferRay(t *testing.T) {
	ctx := gctx.GetInitCtx()
	var usd int64 = 1
	solAmount := getSolInpByUsd(ctx, t, usd)
	tokens := make([]string, 0)
	tokens = append(tokens, "7WX81RVWm7QEbzZhdD1qqrGJhF9aPN5n6bARxCBDiDSF")
	//tokens = append(tokens, ave.AveToken{
	//	Token0Address: "FUAfBo2jgks6gB4Z4LfZkqSZgzNucisEHqnNebaRxM1P",
	//})
	iDex := NewDex(Raydium)
	fee, err := iDex.GetPriorityFee(ctx)
	if err != nil {
		t.Fatal(err)
	}
	pFee := fee.M
	sol := "So11111111111111111111111111111111111111112"
	for i, token := range tokens {
		if i >= 1 {
			break
		}
		quoteReq := model.QuoteBaseReq{
			InputMint:  sol,
			OutputMint: token,
			Amount:     solAmount,
		}
		start := time.Now()
		bps, err := iDex.GetSlippageBps(ctx, quoteReq)
		if err != nil {
			t.Fatal(err)
		}
		input := model.SwapQuoteReq{
			QuoteBaseReq: quoteReq,
			SlippageBps:  bps,
			PriorityFee:  pFee,
		}
		swapRouter, err := iDex.GetSwapRaw(ctx, input, conf.Address)
		if err != nil {
			t.Fatal(err)
		}
		g.Dump(swapRouter)
		start = time.Now()
		// 反序列化交易
		tx, err := solana.TransactionFromBase64(swapRouter.TxRaw)
		if err != nil {
			t.Fatal(err)
		}
		fromAccount, err := GetSolana().GetFrom(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if fromAccount.PublicKey().Equals(key) {
					return fromAccount
				}
				return nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		signature, err := GetSolana().SendTx(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		g.Dump(signature.String())
		g.Dump(time.Since(start))
		out, err := GetSolana().GetTransaction(ctx, signature)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(out.Meta.Status)
		g.Dump(out.Meta.Status["Ok"])
		_, ok := out.Meta.Status["Ok"]
		g.Dump(ok)
	}
}

func TestGetGetSlippageBps(t *testing.T) {
	ctx := gctx.GetInitCtx()
	iDex := NewDex(Jupiter)
	sol := "So11111111111111111111111111111111111111112"
	out := "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"
	bps, err := iDex.GetSlippageBps(ctx, model.QuoteBaseReq{
		InputMint:  sol,
		OutputMint: out,
		Amount:     4218341,
	})
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(bps)
}

func TestGetPrice(t *testing.T) {
	ctx := gctx.GetInitCtx()
	ray := dex.NewRaydium()
	tokens := []string{
		"So11111111111111111111111111111111111111112",
		"6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN",
		"FUAfBo2jgks6gB4Z4LfZkqSZgzNucisEHqnNebaRxM1P",
	}
	for range 5 {
		prices, err := ray.MintPrice(ctx, tokens)
		if err != nil {
			t.Fatal(err)
		}
		g.Dump(prices)
		time.Sleep(1 * time.Second)
	}

}

func getSolInpByUsd(ctx context.Context, t *testing.T, usd int64) uint64 {
	ray := dex.NewRaydium()
	sol := "So11111111111111111111111111111111111111112"
	prices, err := ray.MintPrice(ctx, []string{sol})
	if err != nil {
		t.Fatal(err)
	}
	solNum := decimal.NewFromInt(usd)
	solPrice := decimal.NewFromFloat(prices[sol]).Round(2)
	solDiv := solNum.Div(solPrice).Round(9)
	inp := solDiv.Shift(9)
	return uint64(inp.IntPart())
}

func TestTransferJup(t *testing.T) {
	ctx := gctx.GetInitCtx()
	var usd int64 = 1
	solAmount := getSolInpByUsd(ctx, t, usd)
	//tokens, err := getMemeTokens(ctx, t)
	//if err != nil {
	//	g.Log().Error(ctx, err)
	//	return
	//}
	tokens := make([]string, 0)
	tokens = append(tokens, "2zMMhcVQEXDtdE6vsFS7S7D5oUodfJHE8vd1gnBouauv")
	//tokens = append(tokens, ave.AveToken{
	//	Token0Address: "FUAfBo2jgks6gB4Z4LfZkqSZgzNucisEHqnNebaRxM1P",
	//})
	iDex := NewDex(Jupiter)
	fee, err := iDex.GetPriorityFee(ctx)
	if err != nil {
		t.Fatal(err)
	}
	pFee := fee.M
	sol := "So11111111111111111111111111111111111111112"
	for i, token := range tokens {
		if i >= 1 {
			break
		}
		quoteReq := model.QuoteBaseReq{
			InputMint:  sol,
			OutputMint: token,
			Amount:     solAmount,
		}
		bps, err := iDex.GetSlippageBps(ctx, quoteReq)
		if err != nil {
			t.Fatal(err)
		}
		g.Dump(bps)
		input := model.SwapQuoteReq{
			QuoteBaseReq: quoteReq,
			SlippageBps:  bps,
			PriorityFee:  pFee,
			FeeAccount:   consts.FeeAccountATA,
		}
		swapRouter, err := iDex.GetSwapRaw(ctx, input, conf.Address)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(swapRouter)
		// 反序列化交易
		tx, err := solana.TransactionFromBase64(swapRouter.TxRaw)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		fromAccount, err := GetSolana().GetFrom(ctx)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if fromAccount.PublicKey().Equals(key) {
					return fromAccount
				}
				return nil
			},
		)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		signature, err := GetSolana().SendTx(ctx, tx)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(signature.String())
		out, err := GetSolana().GetTransaction(ctx, signature)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(out.Meta.Status)
		g.Dump(out.Meta.Status["Ok"])
		_, ok := out.Meta.Status["Ok"]
		g.Dump(ok)
	}
}

func TestSellMeme(t *testing.T) {
	ctx := gctx.GetInitCtx()
	accounts, err := GetSolana().GetTokenBalance(ctx, conf.Address)
	if err != nil {
		t.Error(err)
	}
	tokens := make([]string, 0)
	for _, account := range accounts {
		tokens = append(tokens, account.Mint.String())
	}
	//tokens = make([]string, 0)
	//tokens = append(tokens, "CzRYeBnkk2osncaEbw6TkXt2UtBRxmKMcdaZZZxrgoyf")
	iDex := NewDex(Jupiter)
	fee, err := iDex.GetPriorityFee(ctx)
	if err != nil {
		t.Fatal(err)
	}
	pFee := fee.M
	sol := "So11111111111111111111111111111111111111112"
	prices, err := iDex.GetTokenPrice(ctx, strings.Join(tokens, ","))
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(prices)
	for _, account := range accounts {
		if account.Amount == 0 || !garray.NewStrArrayFrom(tokens).Contains(account.Mint.String()) {
			continue
		}
		input := model.SwapQuoteReq{
			QuoteBaseReq: model.QuoteBaseReq{
				InputMint:  account.Mint.String(),
				OutputMint: sol,
				Amount:     account.Amount,
			},
			SlippageBps: 50,
			PriorityFee: pFee,
			FeeAccount:  consts.FeeAccountATA,
		}
		swapRouter, err := iDex.GetSwapRaw(ctx, input, conf.Address)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(swapRouter)
		// 反序列化交易
		tx, err := solana.TransactionFromBase64(swapRouter.TxRaw)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		fromAccount, err := GetSolana().GetFrom(ctx)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if fromAccount.PublicKey().Equals(key) {
					return fromAccount
				}
				return nil
			},
		)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		signature, err := GetSolana().SendTx(ctx, tx)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(signature.String())
		out, err := GetSolana().GetTransaction(ctx, signature)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		g.Dump(out.Meta.Status)
		g.Dump(out.Meta.Status["Ok"])
		_, ok := out.Meta.Status["Ok"]
		g.Dump(ok)
	}
}

func TestCreateMint(t *testing.T) {
	ctx := gctx.GetInitCtx()
	GetSolana().CreateMint(ctx)
}
