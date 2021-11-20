package syncwd

import (
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

type SyncManager struct {
	m  map[string]ModelI
	ml *sync.RWMutex

	pool RedisPoolI

	offsetDay int
}

func NewSyncManager(p *redis.Pool, offsetday int) *SyncManager {
	return &SyncManager{
		m:         make(map[string]ModelI),
		ml:        &sync.RWMutex{},
		pool:      p,
		offsetDay: 3,
	}
}

func NewSyncManagerV2(p RedisPoolI, offsetday int) *SyncManager {
	return &SyncManager{
		m:         make(map[string]ModelI),
		ml:        &sync.RWMutex{},
		pool:      p,
		offsetDay: 3,
	}
}

func (sm *SyncManager) Add(o ModelI) {

	sm.ml.Lock()
	defer sm.ml.Unlock()
	_, exist := sm.m[o.TableName()]
	if exist {
		Errorf("exists job for table name %s\n", o.TableName())
		return
	}

	sm.m[o.TableName()] = o
}

func (sm *SyncManager) Run() {
	sm.ml.RLock()
	var jobs = make([]ModelI, 0, 10)
	for _, v := range sm.m {
		jobs = append(jobs, v)
	}
	sm.ml.RUnlock()

	for _, v := range jobs {
		go sm.handleJob(sm.pool, v)
	}
}

func (sm SyncManager) handleJob(p RedisPoolI, o ModelI) {
	if sm.offsetDay <= 0 {
		sm.offsetDay = 1
	}

	// 对最近的n天的数据，归并进昨天
	merge(p, o, sm.offsetDay)

	// 对昨天单日的数据，进行同步
	handle(p, o)

}

func merge(p RedisPoolI, o ModelI, offsetDays int) error {
	conn := p.Get()
	defer conn.Close()

	latestKey := daySetKey(o, time.Now().AddDate(0, 0, -1))

	var args = make([]interface{}, 0, 10)
	args = append(args, latestKey)

	if offsetDays > 1 {
		for i := 0; i < offsetDays; i ++ {
			key := daySetKey(o, time.Now().AddDate(0, 0, -(i + 1)))
			args = append(args, key)
		}
	}

	//fmt.Printf("merge: %s \n ", args)

	_, e := conn.Do("sunionstore", args...)

	if e != nil {
		return errorx.Wrap(e)
	}

	for _, v := range args {
		if v == latestKey {
			continue
		}
		conn.Send("del", v)
	}

	conn.Flush()

	return nil
}
