package syncwd

import "github.com/garyburd/redigo/redis"

type RedisPoolI interface{
	Get() redis.Conn
}
