package syncwd

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"time"
)

func handle(p *redis.Pool, o ModelI) {
	setkey := daySetKey(o, time.Now().AddDate(0, 0, -1))

    Printf("开始执行 %s 同步计划\n", setkey)

L:
	for {
		done, e := oneloop(p, o, setkey)
		if e != nil {
			fmt.Printf("one loop fail err: %s", errorx.Wrap(e).Error())
		}
		if done {
			break L
		}
	}
}

// 当set里没有记录时，返回true,nil
// 当set有记录，但是一次性可以执行玩，返回true,nil
// 当set仍然有记录时, 返回false,nil
func oneloop(p *redis.Pool, o ModelI, setkey string) (bool, error) {
	conn := p.Get()

	defer conn.Close()

	var script = `
    local total_num = redis.call('scard',KEYS[1]);
    if total_num <= 500 then
        -- 500个以内，直接全部返回
        local rands = redis.call('spop', KEYS[1], total_num);
        return {total_num, rands,0};
    end

    local rands = redis.call('spop', KEYS[1], 500);
    return {total_num, rands, total_num-500};
`
	raw, e := conn.Do("eval", script, 1, setkey)

	rss := raw.([]interface{})

	if e != nil {
		return false, errorx.Wrap(e)
	}

	if len(rss) != 3 {
		return false, errorx.NewFromStringf("bad returned value %v", rss)
	}

	totalLength := rss[0].(int64)
	if totalLength == 0 {
		return true, nil
	}

	keysI := rss[1].([]interface{})
	fmt.Println(Debug(rss[1]))

	for _, vi := range keysI {
		v, _ := redis.String(vi, nil)

		rs, e := redis.Bytes(conn.Do("get", v))
		if e != nil && e == redis.ErrNil {
			continue
		}

		if e != nil {
			return false, errorx.Wrap(e)
		}

		t := reflect.TypeOf(o)
		if t.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
			t = t.Elem()
		}

		instance := reflect.New(t)
		ptr := instance.Interface()

		if e := json.Unmarshal(rs, &ptr); e != nil {
			fmt.Printf("unmarshal err for key %v buf %s \n", v, rs)
			continue
		}

		m := instance.MethodByName("SyncToDB")
		if !m.IsNil() {
			m.Call(nil)
			continue
		}

	}

	restLength := rss[2].(int64)

	if restLength == 0 {
		return true, nil
	}

	return false, nil
}
