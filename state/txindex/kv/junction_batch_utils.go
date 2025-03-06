package kv

import (
	"encoding/binary"
	"fmt"
	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/types"
)

var (
	KeyCountPodsTxs        = []byte("JunctionCountTxs")      // Stores the total transaction count
	KeyCountPods           = []byte("JunctionCountPods")     // Stores the number of transaction batches (pods)
	KeyTxHashesBatchPrefix = []byte("JunctionTxHashesBatch") // Prefix for batch storage
	MaxCountTxBatch        = 128                             // Maximum transactions per batch
)

// InitiateDatabaseForPods 🔹 Initializes the database for storing transaction batches
func InitiateDatabaseForPods(txi *TxIndex) error {

	if txi.store == nil {
		return fmt.Errorf("txi.store is not initialized")
	}

	// Initialize KeyCountPodsTxs (transaction count) to 0
	err := SetCountPodsTxs(txi, 0)
	if err != nil {
		return fmt.Errorf("failed to set KeyCountPodsTxs: %w", err)
	}

	// Initialize KeyCountPods (batch count) to 1 (first batch)
	err = SetCountPods(txi, 1)
	if err != nil {
		return fmt.Errorf("failed to set KeyCountPods: %w", err)
	}

	// Initialize KeyTxHashesBatch (transaction hashes batch) for batch 1
	//err = SetTxHashesBatch(txi, 1, [][]byte{})
	//if err != nil {
	//	return fmt.Errorf("failed to initialize first batch: %w", err)
	//}

	return nil
}

// StoreTxHashesBatch 🔹 Stores transaction hashes in batches of 128
func StoreTxHashesBatch(txi *TxIndex, b *txindex.Batch) error {
	if txi.store == nil {
		return fmt.Errorf("txi.store is not initialized")
	}

	// 🔹 Get current transaction count
	countPodsTxs, err := GetCountPodsTxs(txi)
	if err != nil {
		return fmt.Errorf("failed to get KeyCountPodsTxs: %w", err)
	}

	// 🔹 Get current batch number (countPods)
	countPods, err := GetCountPods(txi)
	if err != nil {
		return fmt.Errorf("failed to get KeyCountPods: %w", err)
	}

	// 🔹 Retrieve existing transaction hashes for the current batch
	existingTxHashes, err := GetTxHashesBatch(txi, countPods)
	if err != nil {
		// First batch initialization (if no batch exists yet)
		existingTxHashes = [][]byte{}
	}

	// 🔹 Process new transaction hashes from b.Ops
	var newTxHashes [][]byte
	for _, result := range b.Ops {
		hashByte := types.Tx(result.Tx).Hash()
		newTxHashes = append(newTxHashes, hashByte)
	}

	// 🔹 Merge existing and new hashes
	allTxHashes := append(existingTxHashes, newTxHashes...)

	// 🔹 Process batches iteratively
	for len(allTxHashes) >= MaxCountTxBatch {
		// ✅ Save a full batch of 128 transactions
		err = SetTxHashesBatch(txi, countPods, allTxHashes[:MaxCountTxBatch])
		if err != nil {
			return fmt.Errorf("failed to store completed batch: %w", err)
		}

		// ✅ Increment batch count
		countPods++
		err = SetCountPods(txi, countPods)
		if err != nil {
			return fmt.Errorf("failed to update KeyCountPods: %w", err)
		}

		// ✅ Remove saved transactions and continue
		allTxHashes = allTxHashes[MaxCountTxBatch:]
	}

	// 🔹 Save remaining transactions in the current batch
	err = SetTxHashesBatch(txi, countPods, allTxHashes)
	if err != nil {
		return fmt.Errorf("failed to store remaining transactions: %w", err)
	}

	// 🔹 Update total transaction count
	countPodsTxs += uint64(len(b.Ops))
	err = SetCountPodsTxs(txi, countPodsTxs)
	if err != nil {
		return fmt.Errorf("failed to update KeyCountPodsTxs: %w", err)
	}

	return nil
}

// SetTxHashesBatch 🔹 Store transaction hashes for a specific batch number
func SetTxHashesBatch(txi *TxIndex, batchNumber uint64, hashes [][]byte) error {
	KeyTxHashesBatch := append(KeyTxHashesBatchPrefix, Uint64ToBytes(batchNumber)...)

	var batchBytes []byte
	for _, hash := range hashes {
		batchBytes = append(batchBytes, hash...)
	}

	return txi.store.Set(KeyTxHashesBatch, batchBytes)
}

// GetTxHashesBatch 🔹 Retrieve transaction hashes for a specific batch number
func GetTxHashesBatch(txi *TxIndex, batchNumber uint64) ([][]byte, error) {
	KeyTxHashesBatch := append(KeyTxHashesBatchPrefix, Uint64ToBytes(batchNumber)...)

	batch, err := txi.store.Get(KeyTxHashesBatch)
	if err != nil || len(batch) == 0 {
		return nil, fmt.Errorf("batch not found")
	}

	var txHashes [][]byte
	for i := 0; i < len(batch); i += 32 { // Assuming 32-byte transaction hashes
		if i+32 <= len(batch) {
			txHashes = append(txHashes, batch[i:i+32])
		}
	}
	return txHashes, nil
}

// Uint64ToBytes 🔹 Convert uint64 to []byte
func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

// BytesToUint64 🔹 Convert []byte to uint64
func BytesToUint64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

// SetCountPodsTxs 🔹 Set KeyCountPodsTxs
func SetCountPodsTxs(txi *TxIndex, count uint64) error {
	return txi.store.Set(KeyCountPodsTxs, Uint64ToBytes(count))
}

// GetCountPodsTxs 🔹 Get total transaction count
func GetCountPodsTxs(txi *TxIndex) (uint64, error) {
	b, err := txi.store.Get(KeyCountPodsTxs)
	if err != nil || len(b) == 0 {
		return 0, err
	}
	return BytesToUint64(b), nil
}

// SetCountPods 🔹 Set KeyCountPods
func SetCountPods(txi *TxIndex, count uint64) error {
	return txi.store.Set(KeyCountPods, Uint64ToBytes(count))
}

// GetCountPods 🔹 Get total batch count
func GetCountPods(txi *TxIndex) (uint64, error) {
	b, err := txi.store.Get(KeyCountPods)
	if err != nil || len(b) == 0 {
		return 0, err
	}
	return BytesToUint64(b), nil
}
