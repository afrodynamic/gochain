package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/afrodynamic/gochain/api/internal/adapter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type Adapter struct {
	rpc *rpcClient
}

func NewAdapter() *Adapter {
	return &Adapter{rpc: newRPC()}
}

func (ad *Adapter) Network() string {
	return "ethereum"
}

func (ad *Adapter) NewKey(seed []byte) (string, string, string, error) {
	seedHash := gethcrypto.Keccak256Hash(seed)
	privateKey, _ := gethcrypto.ToECDSA(seedHash.Bytes())
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	publicKeyBytes := gethcrypto.FromECDSAPub(publicKey)
	addressHex := common.BytesToAddress(gethcrypto.Keccak256(publicKeyBytes[1:])[12:]).Hex()

	return hex.EncodeToString(gethcrypto.FromECDSA(privateKey)), hex.EncodeToString(publicKeyBytes), addressHex, nil
}

func (ad *Adapter) ParseAddress(address string) (string, error) {
	if len(address) == 42 && strings.HasPrefix(strings.ToLower(address), "0x") {
		return address, nil
	}

	return "", errors.New("invalid ethereum address")
}

func (ad *Adapter) Balance(ctx context.Context, address string) (uint64, error) {
	if ad.rpc.url == "" {
		return 0, nil
	}

	var balanceHex string

	if err := ad.rpc.call(ctx, "eth_getBalance", []interface{}{address, "latest"}, &balanceHex); err != nil {
		return 0, err
	}

	balance := new(big.Int)
	balance.SetString(strings.TrimPrefix(balanceHex, "0x"), 16)

	return balance.Uint64(), nil
}

func (ad *Adapter) BuildTx(ctx context.Context, sender, recipient string, amount uint64, feeHint adapter.FeeHint) (adapter.Tx, error) {
	return adapter.Tx{
		From:   sender,
		To:     recipient,
		Amount: amount,
		Fee:    feeHint.MaxFeePerGas,
		Nonce:  0,
		Data:   nil,
	}, nil
}

func (ad *Adapter) SignTx(privateKeyHex string, tx adapter.Tx) (adapter.SignedTx, error) {
	if ad.rpc.url == "" {
		raw := hex.EncodeToString([]byte(tx.From + tx.To))

		return adapter.SignedTx{RawHex: raw, TxID: "0x" + raw[:64]}, nil
	}

	ctx := context.Background()
	var chainIDHex string

	if err := ad.rpc.call(ctx, "eth_chainId", []interface{}{}, &chainIDHex); err != nil {
		return adapter.SignedTx{}, err
	}

	chainID := new(big.Int)
	chainID.SetString(strings.TrimPrefix(chainIDHex, "0x"), 16)

	var nonceHex string

	if err := ad.rpc.call(ctx, "eth_getTransactionCount", []interface{}{tx.From, "pending"}, &nonceHex); err != nil {
		return adapter.SignedTx{}, err
	}

	nonce := new(big.Int)
	nonce.SetString(strings.TrimPrefix(nonceHex, "0x"), 16)

	var tipHex, gasPriceHex string
	_ = ad.rpc.call(ctx, "eth_maxPriorityFeePerGas", []interface{}{}, &tipHex)
	_ = ad.rpc.call(ctx, "eth_gasPrice", []interface{}{}, &gasPriceHex)

	tipWei := new(big.Int)

	if strings.HasPrefix(tipHex, "0x") {
		tipWei.SetString(strings.TrimPrefix(tipHex, "0x"), 16)
	} else {
		tipWei.SetUint64(1_500_000_000)
	}

	baseWei := new(big.Int)

	if strings.HasPrefix(gasPriceHex, "0x") {
		baseWei.SetString(strings.TrimPrefix(gasPriceHex, "0x"), 16)
	} else {
		baseWei.SetUint64(5_000_000_000)
	}

	maxFeePerGas := new(big.Int).Add(baseWei, tipWei)
	toAddress := common.HexToAddress(tx.To)
	valueWei := new(big.Int).SetUint64(tx.Amount)

	dynamicTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce.Uint64(),
		GasTipCap: tipWei,
		GasFeeCap: maxFeePerGas,
		Gas:       21000,
		To:        &toAddress,
		Value:     valueWei,
		Data:      nil,
	}

	txToSign := types.NewTx(dynamicTx)
	privateKeyBytes, _ := hex.DecodeString(strings.TrimPrefix(privateKeyHex, "0x"))
	ecdsaKey, _ := gethcrypto.ToECDSA(privateKeyBytes)
	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(txToSign, signer, ecdsaKey)

	if err != nil {
		return adapter.SignedTx{}, err
	}

	rawBytes, err := signedTx.MarshalBinary()

	if err != nil {
		return adapter.SignedTx{}, err
	}

	return adapter.SignedTx{RawHex: "0x" + hex.EncodeToString(rawBytes), TxID: signedTx.Hash().Hex()}, nil
}

func (ad *Adapter) Broadcast(ctx context.Context, signedTx adapter.SignedTx) (string, error) {
	if ad.rpc.url == "" {
		return signedTx.TxID, nil
	}

	var txID string

	if err := ad.rpc.call(ctx, "eth_sendRawTransaction", []interface{}{signedTx.RawHex}, &txID); err != nil {
		return "", err
	}

	return txID, nil
}

func (ad *Adapter) TxStatus(ctx context.Context, txID string) (adapter.Status, error) {
	if ad.rpc.url == "" {
		return adapter.StatusPending, nil
	}

	type receipt struct {
		BlockNumber string `json:"blockNumber"`
	}

	var rcpt *receipt

	if err := ad.rpc.call(ctx, "eth_getTransactionReceipt", []interface{}{txID}, &rcpt); err != nil {
		return adapter.StatusPending, nil
	}

	if rcpt == nil || rcpt.BlockNumber == "" {
		return adapter.StatusPending, nil
	}

	return adapter.StatusMined, nil
}

var _ adapter.ChainAdapter = (*Adapter)(nil)
