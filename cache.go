package gocache

type Cache interface {
	Set(key string, content []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
	GetStat() Stat
}

type Stat struct {
	KeyCount  uint64
	KeySize   uint64
	ValueSize uint64
}

func (stat *Stat) addKeyStat(key string, value []byte) {
	stat.KeyCount += 1
	stat.KeySize += uint64(len(key))
	stat.ValueSize += uint64(len(value))
}

func (stat *Stat) delKeyStat(key string, value []byte) {
	stat.KeyCount -= 1
	stat.KeySize -= uint64(len(key))
	stat.ValueSize -= uint64(len(value))
}
