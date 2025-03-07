package null

import (
	"context"
	"errors"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/state/txindex"
)

var _ txindex.TxIndexer = (*TxIndex)(nil)

// TxIndex acts as a /dev/null.
type TxIndex struct{}

func (txi *TxIndex) GetBatchArray(batchNumber uint64) ([][]byte, error) {
	var txHashesBatch [][]byte
	return txHashesBatch, errors.New(`indexing is disabled (set 'tx_index = "kv"' in config)`)
}

func (txi *TxIndex) CountPodsTxs() (uint64, error) {
	return 0, errors.New(`indexing is disabled (set 'tx_index = "kv"' in config)`)

}

func (txi *TxIndex) AddPod(b *txindex.Batch) error {
	return nil
}

// Get on a TxIndex is disabled and panics when invoked.
func (txi *TxIndex) Get(hash []byte) (*abci.TxResult, error) {
	return nil, errors.New(`indexing is disabled (set 'tx_index = "kv"' in config)`)
}

// AddBatch is a noop and always returns nil.
func (txi *TxIndex) AddBatch(batch *txindex.Batch) error {
	return nil
}

// Index is a noop and always returns nil.
func (txi *TxIndex) Index(result *abci.TxResult) error {
	return nil
}

func (txi *TxIndex) Search(ctx context.Context, q *query.Query) ([]*abci.TxResult, error) {
	return []*abci.TxResult{}, nil
}
