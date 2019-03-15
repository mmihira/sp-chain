package db

import (
	"bytes"
)

// DbIterface Interface for database access
type Interface interface {
	SaveBlock(blockHash string, buff *bytes.Buffer) error
	// GetBlock(blockhash string) (*bytes.Buffer, error)

	// SaveUTXO(txHash string, buff *bytes.Buffer) error
	// GetUTXO(txHash string) (*bytes.Buffer, error)
}

