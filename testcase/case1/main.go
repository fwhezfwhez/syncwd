package main

import (
	"fmt"
	"github.com/fwhezfwhez/syncwd"
	"github.com/fwhezfwhez/syncwd/testcase/case1/src"
	"github.com/garyburd/redigo/redis"
	"time"
)

var sm *syncwd.SyncManager

var p *redis.Pool

func init() {
	p = src.NewPool("49.234.137.226:6379", "echo123#qp", 0)

	sm = syncwd.NewSyncManager(p, 2)
	sm.Add(src.UserLevelProcess{})
}

func main() {
//	prepareData()
	//

	//c := cron.New()
	//c.AddFunc("0 0 3 * * ?", func() {
	//	sm.Run()
	//})
	//
	//c.Start()

	sm.Run()

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
