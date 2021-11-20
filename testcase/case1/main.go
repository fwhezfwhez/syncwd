package main

import (
	"fmt"
	"github.com/fwhezfwhez/syncwd"
	"github.com/fwhezfwhez/syncwd/testcase/case1/src"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/robfig/cron.v2"
	"time"
)

var sm *syncwd.SyncManager

var p *redis.Pool

func init() {
	p = src.NewPool("localhost:6379", "", 0)

	sm = syncwd.NewSyncManager(p, 5)
	sm.Add(src.UserInfo{})
}

func main() {
	prepareData()
	//

	c := cron.New()
	c.AddFunc("0 0 3 * * ?", func() {
		sm.Run()
	})

	c.Start()

	select {}
}

func prepareData() {
	// 准备数据源
	var sd = syncwd.NewSyncwd()

	conn := p.Get()
	defer conn.Close()
	for i := 0; i < 10000; i ++ {
		ui := src.UserInfo{
			UserName: fmt.Sprintf("冯%d-%d", time.Now().Unix(), i),
		}

		if e := sd.Update(ui, conn); e != nil {
			panic(e)
		}
	}
}
