package leveldb

import (
	"os"
	"testing"
	"spchain/internal"
	"spchain/chain"
	"github.com/davecgh/go-spew/spew"
)

var testDbPath = "/tmp/spchain-test"

func resetDb() {
	// delete file
	os.RemoveAll(testDbPath)
}

func TestLevelDbCreation(t *testing.T) {
	resetDb()

	// createdatabase
	_, err := InitDatabaseAtPath("/tmp/spchain-test")
	if err != nil {
		t.Errorf("Got error attempting to open database %s", err)
	}
}

func TestSaveBlock( t *testing.T) {
	resetDb()

	// createdatabase
	db, err := InitDatabaseAtPath("/tmp/spchain-test")

	if err != nil {
		t.Errorf("Got error attempting to open database %s", err)
	}

	block := internal.BlockWithCoinBase()
	hashStr := block.HashString()

	saveErr := db.SaveBlock(hashStr, block.Ser())
	if saveErr != nil {
		t.Errorf("Got error attempting to save database %s", saveErr)
	}

	getBlock, getError := chain.GetBlock(hashStr, &db)
	if getError != nil {
		t.Errorf("Got error attempting to retrieve  block%s", getError)
	}

	if getBlock.HashString() != block.HashString() {
		t.Errorf("Expected %s block hash got %s", getBlock.HashString(), block.HashString())
	}
}
