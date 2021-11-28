package syncwd

type ModelI interface {
	RedisKey() string
	TableName() string
	SyncToDB() error
}
