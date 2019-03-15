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

// SaveBlock Save a block to the database
func (db LevelDb) SaveBlock(blockhash string, buff *bytes.Buffer) error {
	key := fmt.Sprintf("b_%s", blockhash)
	return db.Db.Put([]byte(key), buff.Bytes(), nil)
}




