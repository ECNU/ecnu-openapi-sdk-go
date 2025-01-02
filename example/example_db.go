package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecnu/ecnu-openapi-sdk-go/sdk"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func exampleSyncToDB() {
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

	// 连接到数据库，可以同步到所有 gorm 支持的数据库
	// 详见 gorm 相关文档：https://gorm.io/zh_CN/docs/connecting_to_the_database.html
	/*
		MySQL Connect
		dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		dial := mysql.Open(dsn)

		Postgres Connect
		dsn := "host=localhost user=user password=pass dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		dial := postgres.Open(dsn)


		SQL Server Connect
		dsn := "sqlserver://user:pass@localhost:1433?database=dbname&encrypt=disable"
		dial := sqlserver.Open(dsn)
	*/
	dial := sqlite.Open("gorm.db")
	db, err := gorm.Open(dial, ormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
		   type APIConfig struct {
		   	APIPath        string `json:"api_path"`   // 接口的地址，例如 /api/v1/organization/list, 也可以追加参数，例如 /api/v1/organization/list?departmentId=0445
		   	PageSize       int    `json:"page_size"`  // 翻页参数会自动添加，默认 pageSize 是 2000，最大值是 10000。
			BatchSize      int    `json:"data_batch"` // 批量写入数据时的批次大小，默认是100。给的太大可能会数据库报错，请根据实际情况调整。
			UpdatedAtField string                     // 增量同步时，数据库内的时间戳字段名，默认是 updated_at
		   }

	*/

	// 配置待同步的接口
	api := sdk.APIConfig{
		APIPath:   "/api/v1/sync/fakewithts",
		PageSize:  2000,
		BatchSize: 100,
	}

	fakeRows := []FakeRowsWithTS{}

	// 如果接口不支持软删除标记，且需要删除上游已删除的数据，可以先删除表，再全量同步
	// 如果希望自己在同步时建立软删除标记，可以建立临时表进行全量同步
	// 再将临时表的数据更新到主表，并对比数据建立软删除标记。
	/*
		if err = db.Migrator().DropTable(rows); err != nil {
			return
		}
	*/

	// 首次同步时，添加参数 ts=0，同步当前全部有效数据
	// 如果未创建表会自动根据 model 建表
	api.SetParam("ts", "0")
	rowsCount, err := sdk.SyncToDB(db, api, &fakeRows)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("DB:首次同步，从接口获取到 %d 条数据\n", rowsCount)

	// 获取数据库内最后一条时间戳
	ts := sdk.GetLastUpdatedTS(db, api, FakeRowsWithTS{})
	fmt.Printf("最后一条时间戳：%d\n", ts)

	// 参照接口文档，添加full参数，获取包含删除的数据
	// 因此可以捕捉上游数据删除的情况，以软删除的形式记录到数据库

	api.SetParam("ts", fmt.Sprintf("%d", ts))
	api.SetParam("full", "1")

	fakeRows = []FakeRowsWithTS{}
	rowsCount, err = sdk.SyncToDB(db, api, &fakeRows)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("DB:增量同步，从接口获取到 %d 条数据\n", rowsCount)

}
