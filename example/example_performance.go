package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/ecnu/ecnu-openapi-sdk-go/sdk"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func exampleCallAPIPerformance(api sdk.APIConfig) {
	startTime := time.Now().UnixMilli()

	fmt.Printf("单次接口调用开始，pageSize=%d\n", api.PageSize)
	apiPath := api.APIPath + "?" + api.ParamEncode()
	c := sdk.GetOpenAPIClient()
	_, err := c.GetRows(apiPath, 1, api.PageSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	endTime := time.Now().UnixMilli()

	fmt.Printf("单次接口结束，pageSize=%d, 用时 %v 秒\n", api.PageSize, float64(endTime-startTime)/1000)
}

func exampleSyncPerformanceByModel(api sdk.APIConfig) {
	fakeRows := []FakeRows{}
	// 创建一个 MemStats 对象
	m := &runtime.MemStats{}
	startTime := time.Now().UnixMilli()
	// 在程序开始前调用 ReadMemStats 函数，获取初始内存信息
	runtime.ReadMemStats(m)
	initMem := m.Alloc // 获取初始分配的内存字节数

	fmt.Printf("Model:首次同步开始\n")
	if err := sdk.SyncToModel(api, &fakeRows); err != nil {
		fmt.Println(err)
		return
	}
	endTime := time.Now().UnixMilli()
	// 在程序结束后调用 ReadMemStats 函数，获取结束内存信息
	runtime.ReadMemStats(m)
	endMem := m.Alloc // 获取结束分配的内存字节数

	// 计算程序运行时分配的内存字节数
	memUsage := endMem - initMem
	fmt.Printf("Model：同步结束，获取 %d 条数据，用时 %v 秒,分配 %v MB内存\n", len(fakeRows), float64(endTime-startTime)/1000, float64(memUsage)/(1024.0*1024.0))
}

func exampleSyncPerformanceByDB(api sdk.APIConfig, db *gorm.DB) {
	//先把表干掉重新初始化
	fakeRows := []FakeRows{}
	if err := db.Migrator().DropTable(&fakeRows); err != nil {
		fmt.Println(err)
		return
	}

	startTime := time.Now().UnixMilli()

	fmt.Printf("%s：首次同步开始\n", db.Name())
	_, err := sdk.SyncToDB(db, api, &fakeRows)
	if err != nil {
		fmt.Println(err)
		return
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("%s：首次同步完成，全量插入用时 %v 秒\n", db.Name(), float64(endTime-startTime)/1000)

	//测试更新性能
	startTime = time.Now().UnixMilli()
	fmt.Printf("%s：第二次同步开始\n", db.Name())
	_, err = sdk.SyncToDB(db, api, &fakeRows)
	if err != nil {
		fmt.Println(err)
		return
	}
	endTime = time.Now().UnixMilli()

	fmt.Printf("%s：第二次同步完成，全量更新用时 %v 秒\n", db.Name(), float64(endTime-startTime)/1000)
}

func exampleSyncPerformance() {
	// 同步到数据库
	// 配置 gorm，详见：https://gorm.io/zh_CN/docs/gorm_config.html
	ormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Warn, // Log level
				Colorful:      true,        // 彩色打印
			},
		),
		SkipDefaultTransaction: true,
	}

	// 配置待同步的接口
	api := sdk.APIConfig{
		APIPath:   "/api/v1/sync/fake",
		PageSize:  10,
		BatchSize: 100,
	}

	api.SetParam("totalNum", "100")
	m := &runtime.MemStats{}
	// 在程序开始前调用 ReadMemStats 函数，获取初始内存信息
	runtime.ReadMemStats(m)
	initMem := m.Alloc // 获取初始分配的内存字节数

	exampleCallAPIPerformance(api)
	exampleSyncPerformanceByModel(api)
	//MySQL Connect
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dial := mysql.Open(dsn)
	dbMySQL, err := gorm.Open(dial, ormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	exampleSyncPerformanceByDB(api, dbMySQL)
	//Postgres Connect
	dsn = "host=localhost user=user password=pass dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dial = postgres.Open(dsn)
	dbPG, err := gorm.Open(dial, ormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	exampleSyncPerformanceByDB(api, dbPG)
	//SQL Server Connect
	dsn = "sqlserver://user:pass@localhost:1433?database=dbname&encrypt=disable"
	dial = sqlserver.Open(dsn)
	dbSQLServer, err := gorm.Open(dial, ormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	exampleSyncPerformanceByDB(api, dbSQLServer)
	dial = sqlite.Open("gorm.db")
	dbSQLlite, err := gorm.Open(dial, ormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	exampleSyncPerformanceByDB(api, dbSQLlite)
	// 在程序结束后调用 ReadMemStats 函数，获取结束内存信息
	runtime.ReadMemStats(m)
	endMem := m.Alloc // 获取结束分配的内存字节数

	// 计算程序运行时分配的内存字节数
	memUsage := endMem - initMem
	fmt.Printf("分配 %v MB内存\n", float64(memUsage)/(1024.0*1024.0))

}
