package syncwd

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Syncwd struct{}

func NewSyncwd() *Syncwd {
	return &Syncwd{}
}

func (s Syncwd) Update(o ModelI, conn redis.Conn) error {
	buf, e := json.Marshal(o)
	if e != nil {
		return errorx.Wrap(e)
	}

	conn.Send("setex", o.RedisKey(), 60*60*24*7, buf)
	conn.Send("sadd", syncDailySetKey(o), o.RedisKey())
	conn.Send("expire", syncDailySetKey(o), 7*60*60*24)

	fmt.Printf("成功将数据 %s 写入 set %s \n", o.RedisKey(), syncDailySetKey(o))
	conn.Flush()


	conn.Receive()
	conn.Receive()
	conn.Receive()


	return nil
}

func Update(o ModelI, conn redis.Conn) error {
	buf, e := json.Marshal(o)
	if e != nil {
		return errorx.Wrap(e)
	}

	conn.Send("setex", o.RedisKey(), 60*60*24*7, buf)
	conn.Send("sadd", syncDailySetKey(o), o.RedisKey())
	conn.Send("expire", syncDailySetKey(o), 7*60*60*24)

	//fmt.Printf("成功将数据 %s 写入 set %s \n", o.RedisKey(), syncDailySetKey(o))
	conn.Flush()

	conn.Receive()
	conn.Receive()
	conn.Receive()


	return nil
}

func Updates(os []ModelI, conn redis.Conn) error {

	var recvnum int

	for _, v := range os {
		buf, e := json.Marshal(v)
		if e != nil {
			return errorx.Wrap(e)
		}

		conn.Send("setex", v.RedisKey(), 60*60*24*7, buf)
		conn.Send("sadd", syncDailySetKey(v), v.RedisKey())
		conn.Send("expire", syncDailySetKey(v), 7*60*60*24)

		recvnum += 3
	}

	//fmt.Printf("成功将数据 %s 写入 set %s \n", o.RedisKey(), syncDailySetKey(o))
	conn.Flush()

	for i := 0; i < recvnum; i++ {
		conn.Receive()
	}
	return nil
}
func ClearDaySet(o ModelI, conn redis.Conn, t time.Time) {
	conn.Do("del", daySetKey(o, t))
}

func syncDailySetKey(o ModelI) string {
	return fmt.Sprintf("syncwd:sync_daily:%s:%s", o.SourceTableName(), time.Now().Format("2006-01-02"))
}

func daySetKey(o ModelI, t time.Time) string {
	return fmt.Sprintf("syncwd:sync_daily:%s:%s", o.SourceTableName(), t.Format("2006-01-02"))
}
