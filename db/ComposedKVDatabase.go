package db

type ComposedKVDatabase struct {
	mem     *MemKVDatabase
	levelDB *LevelDB
}

func NewComposedKVDatabase(filePath string) *ComposedKVDatabase {
	return &ComposedKVDatabase{
		mem:     NewMemKVDatabase(),
		levelDB: NewLevelDB(filePath),
	}
}

func (db *ComposedKVDatabase) Set(key, value []byte) error {
	db.mem.Set(key, value)
	go db.levelDB.Set(key, value)
	return nil
}

func (db *ComposedKVDatabase) Get(key []byte) (value []byte, err error) {
	value, err = db.mem.Get(key)
	if err != nil {
		value, err = db.levelDB.Get(key)
		if err == nil {
			db.mem.Set(key, value)
		}
	}
	return
}

func (db *ComposedKVDatabase) Delete(key []byte) error {
	db.mem.Delete(key)
	return db.levelDB.Delete(key)
}
