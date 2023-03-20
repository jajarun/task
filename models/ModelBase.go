package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"sync"
	"task/redisClient"
	"time"
)

var instanceDb *gorm.DB
var once sync.Once

func getInstanceDb() *gorm.DB {
	once.Do(func() {
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

type ModelBase struct {
	primaryKey   interface{}
	modelCache   bool
	revisionClue string
}

func (u ModelBase) GetPrimaryId() string {
	return "0"
}

func (u ModelBase) IsModelCache() bool {
	return false
}

func (u ModelBase) GetRevisionClue() string {
	return ""
}

func getCache() *redis.Client {
	return redisClient.GetInstanceRedis()
}

func getCacheKey(model interface{}, key string) string {
	t := reflect.TypeOf(model)
	return t.Elem().Name() + "-" + key
}

func FindFirst(model interface{}, conds ...interface{}) {
	db := getInstanceDb()
	db.First(model, conds)
}

func FindByPrimaryKey(model interface{}, primaryKey string) {
	t := reflect.ValueOf(model)
	method := t.MethodByName("IsModelCache")
	isModelCache := false
	if method.IsValid() {
		returnArr := method.Call(nil)
		isModelCache = returnArr[0].Bool()
		fmt.Println("is cache:", returnArr[0].Bool())
	}
	if isModelCache {
		cache := getCache()
		cacheKey := getCacheKey(model, primaryKey)
		fmt.Println("cacheKey:", cacheKey)
		cacheValue, _ := cache.Get(cacheKey).Result()
		if cacheValue == "" {
			fmt.Println("get data from db")
			FindFirst(model, primaryKey)
			b, err := json.Marshal(model)
			if err != nil {
				fmt.Println("json error:", err)
				return
			}
			cacheValue = string(b)
			cache.SetNX(cacheKey, string(b), time.Second*30)
		} else {
			fmt.Println("get data from cache")
			_ = json.Unmarshal([]byte(cacheValue), model)
		}
	} else {
		fmt.Println("get data from db no cache")
		FindFirst(model, primaryKey)
	}

}

func Save(model interface{}) {
	db := getInstanceDb()
	db.Save(model)

	//清除缓存
	t := reflect.ValueOf(model)
	method := t.MethodByName("GetPrimaryId")
	if method.IsValid() {
		returnArr := method.Call(nil)
		primaryKey := returnArr[0].String()
		cache := getCache()
		cacheKey := getCacheKey(model, primaryKey)
		fmt.Println("delete cache key:" + cacheKey)
		cache.Del(cacheKey)
	}
}
