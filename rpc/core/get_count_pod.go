package core

import (
	"fmt"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/state/txindex/null"
)

func GetCountPods(ctx *rpctypes.Context) (*ctypes.ResultPodCount, error) {
	if _, ok := env.TxIndexer.(*null.TxIndex); ok {
		return nil, fmt.Errorf("transaction indexing is disabled")
	}
	podCount, err := env.TxIndexer.CountPodsTxs()
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultPodCount{
		Count: podCount,
	}, nil
}
