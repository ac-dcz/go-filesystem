package store

type KeyValue interface {
	Hash() string
}

type Store[T KeyValue] interface {
	Get(key string) (T, bool)
	Put(key string, val T)
	Keys() []string
}
