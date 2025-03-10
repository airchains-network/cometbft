package core

import (
	"encoding/hex"
	"fmt"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/state/txindex/null"
)

func GetBatch(ctx *rpctypes.Context, batchNumber uint64) (*ctypes.ResultGetBatch, error) {

	// if index is disabled, return error
	if _, ok := env.TxIndexer.(*null.TxIndex); ok {
		return nil, fmt.Errorf("transaction indexing is disabled")
	}

	r, err := env.TxIndexer.GetBatchArray(batchNumber)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("batch (%X) not found", batchNumber)
	}

	var txs []*ctypes.ResultTx
	var txHashes []string
	for _, hash := range r {

		tx, err := Tx(ctx, hash, true)
		if err != nil {
			return nil, err
		} else {
			txs = append(txs, tx)
			hashStr := hex.EncodeToString(hash)
			txHashes = append(txHashes, hashStr)
		}
	}

	return &ctypes.ResultGetBatch{
		TxCount:      0,
		Transactions: txs,
		TxHashes:     txHashes,
	}, nil

}
