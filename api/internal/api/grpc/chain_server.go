package grpcapi

import (
	"context"

	"github.com/afrodynamic/gochain/api/internal/core"
	chainv1 "github.com/afrodynamic/gochain/api/proto/chain/v1"
)

type ChainServer struct {
	chainv1.UnimplementedChainServer
	blockchain core.Blockchain
}

func NewChain(blockchain core.Blockchain) *ChainServer {
	return &ChainServer{blockchain: blockchain}
}

func (server *ChainServer) GetBlock(ctx context.Context, request *chainv1.GetBlockRequest) (*chainv1.GetBlockResponse, error) {
	block, err := server.blockchain.GetBlock(request.Height)

	if err != nil {
		return nil, err
	}

	return &chainv1.GetBlockResponse{
		Hash:     block.Hash,
		Height:   block.Height,
		PrevHash: block.PrevHash,
	}, nil
}

func (server *ChainServer) SubmitTx(ctx context.Context, request *chainv1.SubmitTxRequest) (*chainv1.SubmitTxResponse, error) {
	txHash, err := server.blockchain.SubmitTx(core.Tx{
		From:   request.From,
		To:     request.To,
		Amount: request.Amount,
		Data:   request.Data,
	})

	if err != nil {
		return nil, err
	}

	return &chainv1.SubmitTxResponse{TxHash: txHash}, nil
}

func (server *ChainServer) GetBalance(ctx context.Context, request *chainv1.GetBalanceRequest) (*chainv1.GetBalanceResponse, error) {
	balance, err := server.blockchain.GetBalance(request.Address)

	if err != nil {
		return nil, err
	}

	return &chainv1.GetBalanceResponse{Balance: balance}, nil
}
