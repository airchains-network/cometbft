package core

import (
	"encoding/json"
	"fmt"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/state/txindex/null"
)

func GetTxHashesByBatch(ctx *rpctypes.Context, batchNumber uint64) (*ctypes.ResultGetTxHashesByBatch, error) {

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

	var txHashes []string
	var txCount uint64
	txCount = 0
	for _, hash := range r {
		tx, err := Tx(ctx, hash, true)
		if err != nil {
			return nil, err
		}
		if tx == nil {
			return nil, fmt.Errorf("transaction not found")
		}

		log := tx.TxResult.GetLog()
		txHash, err := getEthereumTxHash(log)
		if err != nil {
			return nil, err
		}

		txHashes = append(txHashes, txHash)
		txCount++
	}

	return &ctypes.ResultGetTxHashesByBatch{
		TxCount:  txCount,
		TxHashes: txHashes,
	}, nil

}

type TxLog []struct {
	MsgIndex int `json:"msg_index"`
	Events   []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"attributes"`
	} `json:"events"`
}

func getEthereumTxHash(logStr string) (string, error) {
	var txLog TxLog
	err := json.Unmarshal([]byte(logStr), &txLog)
	if err != nil {
		return "", err
	}
	for _, entry := range txLog {
		for _, event := range entry.Events {
			if event.Type == "ethereum_tx" {
				for _, attribute := range event.Attributes {
					if string(attribute.Key) == "ethereumTxHash" {
						return string(attribute.Value), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("ethereumTxHash not found")
}
