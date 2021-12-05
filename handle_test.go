package syncwd

import (
	"fmt"
	"testing"
)

type UserInfo struct {
	UserName string `gorm:"column:user_name;default:"`
}

func (ui UserInfo) SourceTableName() string {
	return "user_info"
}
func (ui UserInfo) RedisKey() string {
	return fmt.Sprintf("testui:%s", ui.UserName)
}

func (ui UserInfo) SyncToDB() error {
	fmt.Println("成功同步:", Debug(ui))
	return nil
}



func TestOneloop(t *testing.T) {
	p := NewPool("localhost:6379", "", 0)

	done, e := oneloop(p, UserInfo{}, "syncwd:sync_daily:user_info:2021-11-20")
	if e != nil {
		panic(e)
	}
	fmt.Println(done)
}

func TestMerge(t *testing.T) {
	p := NewPool("localhost:6379", "", 0)
	if e := merge(p, UserInfo{}, 5); e != nil {
		panic(e)
	}
}
