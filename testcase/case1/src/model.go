package src

import (
	"encoding/json"
	"fmt"
)

type UserInfo struct {
	UserName string `gorm:"column:user_name;default:"`
}

func (ui UserInfo) SourceTableName() string {
	return "user_info"
}

func (ui UserInfo) SyncToDB() error {
	rs, _:=json.MarshalIndent(ui, "  ","  ")
	fmt.Println("成功同步:", string(rs))
	return nil
}

func (ui UserInfo) RedisKey() string {
	return fmt.Sprintf("testui:%s", ui.UserName)
}
