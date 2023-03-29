package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
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
		//_ = db.Callback().Query().Register("after_query", AfterQuery)
		//_ = db.Callback().Update().Register("after_save", AfterSave)
		//_ = db.Callback().Create().Register("after_create", AfterCreate)
		//_ = db.Callback().Delete().Register("after_delete", AfterDelete)
		instanceDb = db
	})
	return instanceDb
}

func AfterQuery(db *gorm.DB) {

}

func flushSingleModel(model ModelInterface) {

}

func AfterCreate(db *gorm.DB) {
	fmt.Println("after create")
}

func AfterSave(db *gorm.DB) {
	fmt.Println("after save")
	flushCache(db)
}

func AfterDelete(db *gorm.DB) {
	fmt.Println("after delete")
	flushCache(db)
}

func flushCache(db *gorm.DB) {
	fmt.Printf("models %v \n", db.Statement.Model)
	return

	var models []interface{}
	switch db.Statement.Model.(type) {
	case []interface{}:
		models = db.Statement.Model.([]interface{})
		break
	case interface{}:
		models = []interface{}{db.Statement.Model.(interface{})}
		break
	default:
		panic("err ty")
	}
	for _, value := range models {
		model := value.(ModelInterface)
		primaryValue := getPrimaryId(model)

		cache := getCache()
		cacheKey := getCacheKey(model, primaryValue)
		fmt.Println("delete cache key:" + cacheKey)
		cache.Del(cacheKey)
	}
	//fmt.Println("delete cache key:" + cacheKey)

	//model := db.Statement.Model.(ModelInterface)
	//primaryValue := getPrimaryId(model)
	//
	//cache := getCache()
	//cacheKey := getCacheKey(model, primaryValue)
	//fmt.Println("delete cache key:" + cacheKey)
	//cache.Del(cacheKey)
}

func getPrimaryId(model ModelInterface) string {
	fmt.Println(model)
	primaryKey := model.getPrimaryKey()
	primaryValue := ""
	r := reflect.ValueOf(model).Elem().FieldByName(primaryKey)
	kType := r.Kind().String()
	switch kType {
	case "string":
		primaryValue = r.String()
		break
	case "int":
	case "uint":
		primaryValue = fmt.Sprintf("%d", r.Uint())
		break
	default:
		panic("primary type error")
	}
	return primaryValue
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
	case int:
		return fmt.Sprintf("%s-%d", modelName, cond.(int))
	case string:
		return fmt.Sprintf("%s-%s", modelName, cond.(string))
	default:
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
