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
	txCount, err := env.TxIndexer.CountPodsTxs()
	if err != nil {
		return nil, err
	}

	// 1 pod can have max 128 transactions.
	maxTxPerPod := uint64(128)
	podCount := countPods(txCount, maxTxPerPod)

	return &ctypes.ResultPodCount{
		TxCount:  txCount,
		PodCount: podCount,
	}, nil
}
func countPods(txCount, maxTxPerPod uint64) uint64 {
	if txCount == 0 || maxTxPerPod == 0 {
		return 0
	}

	pods := txCount / maxTxPerPod
	if txCount%maxTxPerPod != 0 {
		pods++
	}
	return pods
}
