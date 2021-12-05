package syncwd

type ModelI interface {
	RedisKey() string
	SourceTableName() string
	SyncToDB() error
}
