package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"sync"
	"task/components"
	"time"
)

var instanceDb *gorm.DB
var onceDb sync.Once

func GetInstanceDb() *gorm.DB {
	onceDb.Do(func() {
		username := viper.Get("mysql.username") // 账号
		password := viper.Get("mysql.password") // 密码
		host := viper.Get("mysql.host")         // 地址
		port := viper.Get("mysql.port")         // 端口
		DBname := viper.Get("mysql.DBname")     // 数据库名称
		timeout := "10s"                        // 连接超时，10秒
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
		// Open 连接
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		fmt.Println("connect db success")
		instanceDb = db
	})
	return instanceDb
}

func FlushCache(model ModelInterface) {
	primaryKey := model.getPrimaryKey()
	b, _ := json.Marshal(model)
	data := make(map[string]interface{})
	_ = json.Unmarshal(b, &data)
	primaryValue := data[primaryKey]
	cache := getCache()
	cacheKey := getCacheKey(model, primaryValue)
	fmt.Println("delete cache key:" + cacheKey)
	cache.Del(cacheKey)
}

type ModelInterface interface {
	getPrimaryKey() string   //获取主键名
	isModelCache() bool      //是否开启缓存
	getRevisionClue() string //
}

type ModelBase struct {
}

func (mb *ModelBase) getPrimaryKey() string {
	return "ID"
}

func (mb *ModelBase) isModelCache() bool {
	return false
}

func (mb *ModelBase) getRevisionClue() string {
	return ""
}

func getCache() *redis.Client {
	return components.GetInstanceRedis()
}

func getCacheKey(model ModelInterface, cond interface{}) string {
	modelName := reflect.TypeOf(model).Elem().Name()
	switch cond.(type) {
	case int, int32, int64, uint:
		return fmt.Sprintf("%s-%d", modelName, cond.(int))
	case float64:
		return fmt.Sprintf("%s-%d", modelName, int(cond.(float64)))
	case string:
		return fmt.Sprintf("%s-%s", modelName, cond.(string))
	default:
		fmt.Println(reflect.TypeOf(cond))
		panic("condition err")
	}
}

func FindFirst(model ModelInterface, conds ...interface{}) {
	db := GetInstanceDb()
	db.First(model, conds)
}

func FindFirstViaCache(model ModelInterface, conds interface{}, revisionClue string) {

}

func Create(model interface{}) {
	db := GetInstanceDb()
	db.Create(model)
}

func Find(models interface{}, conds ...interface{}) {
	db := GetInstanceDb()
	db.Find(models, conds)
}

func FindByPrimaryKey(model ModelInterface, primaryKey string) {
	isModelCache := model.isModelCache()
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

func Save(model ModelInterface) {
	db := GetInstanceDb()
	db.Save(model)
}

func AutoMigrate(model interface{}) {
	db := GetInstanceDb()
	db.AutoMigrate(model)
}
