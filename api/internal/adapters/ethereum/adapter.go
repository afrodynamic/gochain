package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type Adapter struct{ rpc *rpcClient }

func NewAdapter() *Adapter { return &Adapter{rpc: newRPC()} }

func (a *Adapter) Network() string { return "ethereum" }

func (a *Adapter) NewKey(seed []byte) (string, string, string, error) {
	h := gethcrypto.Keccak256Hash(seed)
	priv, _ := gethcrypto.ToECDSA(h.Bytes())
	pub := priv.Public().(*ecdsa.PublicKey)
	pubBytes := gethcrypto.FromECDSAPub(pub)
	addrHex := common.BytesToAddress(gethcrypto.Keccak256(pubBytes[1:])[12:]).Hex()
	return hex.EncodeToString(gethcrypto.FromECDSA(priv)), hex.EncodeToString(pubBytes), addrHex, nil
}

func (a *Adapter) ParseAddress(s string) (string, error) {
	if len(s) == 42 && strings.HasPrefix(strings.ToLower(s), "0x") {
		return s, nil
	}
	return "", errors.New("invalid ethereum address")
}

func (a *Adapter) Balance(ctx context.Context, addr string) (uint64, error) {
	if a.rpc.url == "" {
		return 0, nil
	}
	var hexBal string
	if err := a.rpc.call(ctx, "eth_getBalance", []interface{}{addr, "latest"}, &hexBal); err != nil {
		return 0, err
	}
	bi := new(big.Int)
	bi.SetString(strings.TrimPrefix(hexBal, "0x"), 16)
	return bi.Uint64(), nil
}

func (a *Adapter) BuildTx(ctx context.Context, from, to string, amt uint64, feeHint core.FeeHint) (core.Tx, error) {
	return core.Tx{From: from, To: to, Amount: amt, Fee: feeHint.MaxFeePerGas, Nonce: 0, Data: nil}, nil
}

func (a *Adapter) SignTx(priv string, tx core.Tx) (core.SignedTx, error) {
	if a.rpc.url == "" {
		raw := hex.EncodeToString([]byte(tx.From + tx.To))
		return core.SignedTx{RawHex: raw, TxID: "0x" + raw[:64]}, nil
	}
	ctx := context.Background()
	var chainHex string
	if err := a.rpc.call(ctx, "eth_chainId", []interface{}{}, &chainHex); err != nil {
		return core.SignedTx{}, err
	}
	chainID := new(big.Int)
	chainID.SetString(strings.TrimPrefix(chainHex, "0x"), 16)

	var nonceHex string
	if err := a.rpc.call(ctx, "eth_getTransactionCount", []interface{}{tx.From, "pending"}, &nonceHex); err != nil {
		return core.SignedTx{}, err
	}
	nonce := new(big.Int)
	nonce.SetString(strings.TrimPrefix(nonceHex, "0x"), 16)

	var tipHex, gasPriceHex string
	_ = a.rpc.call(ctx, "eth_maxPriorityFeePerGas", []interface{}{}, &tipHex)
	_ = a.rpc.call(ctx, "eth_gasPrice", []interface{}{}, &gasPriceHex)

	tip := new(big.Int)
	if strings.HasPrefix(tipHex, "0x") {
		tip.SetString(strings.TrimPrefix(tipHex, "0x"), 16)
	} else {
		tip.SetUint64(1_500_000_000)
	}
	base := new(big.Int)
	if strings.HasPrefix(gasPriceHex, "0x") {
		base.SetString(strings.TrimPrefix(gasPriceHex, "0x"), 16)
	} else {
		base.SetUint64(5_000_000_000)
	}
	maxFee := new(big.Int).Add(base, tip)

	toAddr := common.HexToAddress(tx.To)
	value := new(big.Int).SetUint64(tx.Amount)
	dyn := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce.Uint64(),
		GasTipCap: tip,
		GasFeeCap: maxFee,
		Gas:       21000,
		To:        &toAddr,
		Value:     value,
		Data:      nil,
	}
	gtx := types.NewTx(dyn)
	keyBytes, _ := hex.DecodeString(strings.TrimPrefix(priv, "0x"))
	key, _ := gethcrypto.ToECDSA(keyBytes)
	signer := types.LatestSignerForChainID(chainID)
	signed, err := types.SignTx(gtx, signer, key)
	if err != nil {
		return core.SignedTx{}, err
	}
	rawBytes, err := signed.MarshalBinary()
	if err != nil {
		return core.SignedTx{}, err
	}
	return core.SignedTx{RawHex: "0x" + hex.EncodeToString(rawBytes), TxID: signed.Hash().Hex()}, nil
}

func (a *Adapter) Broadcast(ctx context.Context, stx core.SignedTx) (string, error) {
	if a.rpc.url == "" {
		return stx.TxID, nil
	}
	var txid string
	if err := a.rpc.call(ctx, "eth_sendRawTransaction", []interface{}{stx.RawHex}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

func (a *Adapter) TxStatus(ctx context.Context, id string) (core.Status, error) {
	if a.rpc.url == "" {
		return core.StatusPending, nil
	}
	type receipt struct {
		BlockNumber string `json:"blockNumber"`
	}
	var r *receipt
	if err := a.rpc.call(ctx, "eth_getTransactionReceipt", []interface{}{id}, &r); err != nil {
		return core.StatusPending, nil
	}
	if r == nil || r.BlockNumber == "" {
		return core.StatusPending, nil
	}
	return core.StatusMined, nil
}
