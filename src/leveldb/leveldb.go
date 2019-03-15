package leveldb
import (
	"bytes"
	"fmt"
  "github.com/syndtr/goleveldb/leveldb"
)

type LevelDb struct {
	Db *leveldb.DB
}

// InitDatabaseAtPath
func InitDatabaseAtPath(path string) (LevelDb, error) {
  db, err := leveldb.OpenFile(path, nil)

	if err != nil {
		return LevelDb{}, err
	}
	return LevelDb { Db: db, }, nil
}

func blockId(blockhash string) string {
	return fmt.Sprintf("b_%s", blockhash)
}

// SaveBlock Save a block to the database
func (db LevelDb) SaveBlock(blockhash string, buff *bytes.Buffer) error {
	return db.Db.Put([]byte(blockId(blockhash)), buff.Bytes(), nil)
}

// Get a Block
func (db LevelDb) GetBlock(blockHash string) (*bytes.Buffer, error) {
  data, err := db.Db.Get([]byte(blockId(blockHash)), nil)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	return bytes.NewBuffer(data), nil
}




