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
		_ = db.Callback().Update().Register("after:update", afterSave)
		_ = db.Callback().Delete().Register("after:delete", afterSave)
		//_ = db.Callback().Create().Register("after:create", afterSave)
		fmt.Println("connect db success")
		instanceDb = db
	})
	return instanceDb
}

func afterSave(db *gorm.DB) {
	models := db.Statement.Model
	rm := reflect.ValueOf(models).Elem()
	if rm.Kind() == reflect.Slice || rm.Kind() == reflect.Array {
		length := rm.Len()
		for i := 0; i < length; i++ {
			model := rm.Index(i)
			m := model.Addr().Interface().(ModelInterface)
			FlushCache(m)
		}
	} else {
		m := rm.Addr().Interface().(ModelInterface)
		FlushCache(m)
	}
}

func FlushCache(model ModelInterface) {
	primaryKey := model.GetPrimaryKey()
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
	GetPrimaryKey() string   //获取主键名
	IsModelCache() bool      //是否开启缓存
	GetRevisionClue() string //
}

type ModelBase struct {
}

func (mb *ModelBase) GetPrimaryKey() string {
	return "ID"
}

func (mb *ModelBase) IsModelCache() bool {
	return false
}

func (mb *ModelBase) GetRevisionClue() string {
	return ""
}

func getCache() *redis.Client {
	return components.GetInstanceRedis()
}

func getCacheKey(model ModelInterface, cond interface{}) string {
	modelName := ""
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		modelName = reflect.TypeOf(model).Elem().Name()
	} else {
		modelName = reflect.TypeOf(model).Name()
	}
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
	isModelCache := model.IsModelCache()
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
