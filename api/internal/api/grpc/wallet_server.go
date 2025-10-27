package grpcapi

import (
	"context"
	"time"

	"github.com/afrodynamic/gochain/api/internal/adapter"
	walletv1 "github.com/afrodynamic/gochain/api/proto/wallet/v1"
)

type WalletServer struct {
	walletv1.UnimplementedWalletServer
	adapter adapter.ChainAdapter
}

func NewWallet(ad adapter.ChainAdapter) *WalletServer {
	return &WalletServer{adapter: ad}
}

func NewWalletDummy() *WalletServer {
	return &WalletServer{adapter: nil}
}

func (server *WalletServer) NewKey(ctx context.Context, req *walletv1.NewKeyRequest) (*walletv1.NewKeyResponse, error) {
	if server.adapter == nil {
		return &walletv1.NewKeyResponse{Priv: "priv", Pub: "pub", Addr: "addr"}, nil
	}

	privateKey, publicKey, address, err := server.adapter.NewKey(req.Seed)

	if err != nil {
		return nil, err
	}

	return &walletv1.NewKeyResponse{Priv: privateKey, Pub: publicKey, Addr: address}, nil
}

func (server *WalletServer) ParseAddress(ctx context.Context, request *walletv1.ParseAddressRequest) (*walletv1.ParseAddressResponse, error) {
	if server.adapter == nil {
		return &walletv1.ParseAddressResponse{Addr: request.Value}, nil
	}

	address, err := server.adapter.ParseAddress(request.Value)

	if err != nil {
		return nil, err
	}

	return &walletv1.ParseAddressResponse{Addr: address}, nil
}

func (server *WalletServer) Balance(ctx context.Context, req *walletv1.BalanceRequest) (*walletv1.BalanceResponse, error) {
	if server.adapter == nil {
		return &walletv1.BalanceResponse{Balance: 0}, nil
	}

	balance, err := server.adapter.Balance(ctx, req.Addr)

	if err != nil {
		return nil, err
	}

	return &walletv1.BalanceResponse{Balance: balance}, nil
}

func (server *WalletServer) BuildTx(ctx context.Context, request *walletv1.BuildTxRequest) (*walletv1.BuildTxResponse, error) {
	if server.adapter == nil {
		return &walletv1.BuildTxResponse{Tx: &walletv1.Tx{From: request.From, To: request.To, Amount: request.Amount}}, nil
	}

	tx, err := server.adapter.BuildTx(ctx, request.From, request.To, request.Amount, adapter.FeeHint{
		MaxFeePerGas:   request.FeeHint.MaxFeePerGas,
		MaxPriorityFee: request.FeeHint.MaxPriorityFee,
	})

	if err != nil {
		return nil, err
	}

	return &walletv1.BuildTxResponse{
		Tx: &walletv1.Tx{
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Amount,
			Fee:    tx.Fee,
			Nonce:  tx.Nonce,
			Data:   tx.Data,
		},
	}, nil
}

func (server *WalletServer) SignTx(ctx context.Context, request *walletv1.SignTxRequest) (*walletv1.SignTxResponse, error) {
	if server.adapter == nil {
		return &walletv1.SignTxResponse{Signed: &walletv1.SignedTx{RawHex: "0x", TxId: "id"}}, nil
	}

	signed, err := server.adapter.SignTx(request.Priv, adapter.Tx{
		From:   request.Tx.From,
		To:     request.Tx.To,
		Amount: request.Tx.Amount,
		Fee:    request.Tx.Fee,
		Nonce:  request.Tx.Nonce,
		Data:   request.Tx.Data,
	})

	if err != nil {
		return nil, err
	}

	return &walletv1.SignTxResponse{
		Signed: &walletv1.SignedTx{
			RawHex: signed.RawHex,
			TxId:   signed.TxID,
		},
	}, nil
}

func (server *WalletServer) Broadcast(ctx context.Context, request *walletv1.BroadcastRequest) (*walletv1.BroadcastResponse, error) {
	if server.adapter == nil {
		return &walletv1.BroadcastResponse{TxId: request.Signed.TxId}, nil
	}

	txID, err := server.adapter.Broadcast(ctx, adapter.SignedTx{
		RawHex: request.Signed.RawHex,
		TxID:   request.Signed.TxId,
	})

	if err != nil {
		return nil, err
	}

	return &walletv1.BroadcastResponse{TxId: txID}, nil
}

func (server *WalletServer) TxStatus(ctx context.Context, request *walletv1.TxStatusRequest) (*walletv1.TxStatusResponse, error) {
	if server.adapter == nil {
		return &walletv1.TxStatusResponse{Status: "pending"}, nil
	}

	status, err := server.adapter.TxStatus(ctx, request.TxId)

	if err != nil {
		return nil, err
	}

	return &walletv1.TxStatusResponse{Status: string(status)}, nil
}

func (server *WalletServer) SubscribeTx(request *walletv1.SubscribeTxRequest, stream walletv1.Wallet_SubscribeTxServer) error {
	context := stream.Context()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	transactionID := request.GetId()
	var lastStatus string

	for {
		select {
		case <-context.Done():
			return context.Err()

		case <-ticker.C:
			if server.adapter == nil {
				_ = stream.Send(&walletv1.TxEvent{Id: transactionID, Status: "mined"})
				return nil
			}

			status, err := server.adapter.TxStatus(context, transactionID)

			if err != nil {
				continue
			}

			currentStatus := string(status)

			if currentStatus != lastStatus {
				lastStatus = currentStatus
				_ = stream.Send(&walletv1.TxEvent{Id: transactionID, Status: lastStatus})

				if lastStatus == "mined" {
					return nil
				}
			}
		}
	}
}
