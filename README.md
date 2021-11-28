syncwd

syncwd是一款redis异步落库框架。

## 前言
为了解决消息推送频率过高，db写入数据过载的问题。决定将业务内，高频更新的相关进度，redis化，并组织在第二天低谷凌晨2点开始执行同步。

在使用mq削峰后，仍旧存在极高的写入qps。

## 分析
● 【业务实时性】从实时上考虑，异步写入方式，使用mq削峰后，仍旧很高qps，继续降低拉取频率，只会降低实时性。所以考虑将可靠数据，从db转移至redis。
● 【降幂】异步存储的场景，大部分可以归结为某一个活动进度。这类模型有一个特性，就是以最新的结果为主，换句话说，不建议将每次更新，以日志的形式写入mq，再对每条update分次处理。一个用户单日更新一万次，实际上更新落盘只在乎最新的一次，而不是也在低峰值下，同步一万次。
● 【方案易用性和复制性】 该方案的实现，必须通用化，不能每一个业务场景，每一个业务模块，都重复实现。因为实现过程容易出错，新人和不同的开发同事，实现风格迥异，难以维护。

## 设计方案
### 主流程分析

● 通过将可靠模型存储，从db转移到redis。意味着模型的缓存失效时间，从原来的5-12分钟，升级为3-7天。保障【实时性】。
● 通过redis set集合特性，执行时，用户的进度key，实际上是去重了的，所以能够做到，只取最新一次。保障了【降幂】
● 将复杂的同步过程，形成框架，在集成进项目模块时，越简单，维护起来越简单。保障【方案易用性】


## 实现
● 仓库: github.com/fwhezfwhez/syncwd
● 通过开源，由社区反馈意见和bug。并在团队内，使用优化版的。意见可以通过github issue投递，也可以直接通过邮件1728565484@qq.com投递给作者。

接入
纳入异步同步的表模型，必须实现以下方法
type ModelI interface {
	RedisKey() string            // 某个模型的rediskey
	TableName() string           // 某个模型的表名
	SyncToDB() error             // 执行更新进数据库的方法
}

样例
● 用户活动进度表需要异步，优先更新redis，晚上异步落盘
type UserProcess struct {
	UserName string `gorm:"column:user_name;default:"`
}

func (up UserProcess) TableName() string {
	return "user_process"
}
func (up UserProcess) RedisKey() string {
	return fmt.Sprintf("appsrv:%s", up.UserName)
}

func (up UserProcess) SyncToDB() error {
    // 伪代码，同步过程。实际应该由开发人员自己实现
	fmt.Println("成功同步:", Debug(up))
	return nil
}

func Debug(v interface{}) string {
	rs, _ := json.MarshalIndent(v, "  ", "  ")
	return string(rs)
}
接入同步计划
● 实际测试时，需要先执行prepareData()，再修改本地时间至明天，然后执行
package main

import (
	"fmt"
	"github.com/fwhezfwhez/syncwd"
	"github.com/fwhezfwhez/syncwd/testcase/case1/src"
	"github.com/garyburd/redigo/redis"
	 "gopkg.in/robfig/cron.v2"
	"time"
)

// 同步计划管理 应用内单例。支持并发(cron内的同步操作毋须做任务幂等)。
var sm *syncwd.SyncManager

var p *redis.Pool

func init() {
	p = src.NewPool("localhost:6379", "", 0)

	sm = syncwd.NewSyncManager(p, 5)
	sm.Add(src.UserInfo{})
}

func main() {
	//prepareData()

	c := cron.New()
	c.AddFunc("0 0 3 * * ?", func() {
		sm.Run()
	})

	c.Start()
	select {}
}

// 实际业务中，不需要用到它
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
