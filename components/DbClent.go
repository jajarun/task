package components

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var instanceDb *gorm.DB
var onceDb sync.Once

func GetInstanceDb() *gorm.DB {
	onceDb.Do(func() {
		username := "root"  // 账号
		password := ""      // 密码
		host := "127.0.0.1" // 地址
		port := 3306        // 端口
		DBname := "test"    // 数据库名称
		timeout := "10s"    // 连接超时，10秒
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
		// Open 连接
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		fmt.Println("connect db success")
		if err != nil {
			panic("failed to connect mysql.")
		}
		instanceDb = db
	})
	return instanceDb
}
