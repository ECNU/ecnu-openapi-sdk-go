# ECNU-OpenAPI-Go-SDK

## 能力
- 授权模式（token 管理）
  - [x] client_credentials 模式
  - [ ] password 模式
  - [ ] authorization code 模式
- 接口调用
  - [x] GET 
  - [ ] POST
  - [ ] PUT
  - [ ] DELETE
- 数据同步（接口必须支持翻页）
  - 全量同步
    - [x] 同步为 csv 格式
    - [ ] 同步为 xls/xlsx 格式
    - [x] 同步到数据库
    - [x] 同步到模型
  - 增量同步（接口必须支持ts增量参数）
    - [x] 同步到数据库
    - [x] 同步到模型 

## 依赖
- Go 1.20+
- gorm 1.25+

## 相关资料
- [oauth2.0](https://oauth.net/2/)
- [gorm](https://gorm.io/zh_CN/docs/index.html)

## 支持的数据库
理论上只要 gorm 支持的数据库驱动——[GORM:连接到数据库](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)

都可以支持，以下是测试的情况

如果 gorm 无法直接支持，可以先同步到模型，然后自行处理数据入库的逻辑。

| 数据库 | 驱动 | 测试情况 |
| --- | --- | --- |
| MySQL | gorm.io/driver/mysql | 测试通过 |
| SQLite | github.com/glebarez/sqlite | 测试通过 |
| PostgreSQL | gorm.io/driver/postgres | 测试通过 |
| SQL Server | gorm.io/driver/sqlserver | 测试通过 |

## 示例

### authorization code

Todo

### client_credentials
#### 接口调用
初始化 SDK 后直接调用接口即可，sdk 会自动接管 token 的有效期和续约管理。

```golang
	cf := sdk.OAuth2Config{
		ClientId:     "client_id",
		ClientSecret: "client_secret",
	}
	sdk.InitOAuth2ClientCredentials(cf)
	// 直接调用接口
	res, err := sdk.CallAPI("https://api.ecnu.edu.cn/api/v1/sync/fakewithts?pageNum=1&pageSize=5&ts=0", "GET", nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))
```

#### 数据同步
只需要定义好 orm 映射，SDK 会接管接口调用，数据表创建，数据同步等所有工作。

```golang

	type FakeRowsWithTS struct {
		Id          int       `json:"id" gorm:"primarykey;autoIncrement:false"`
		CreateTime  time.Time `json:"created_at" time_format:"sql_datetime" time_location:"shanghai"`
		UpdateTime  time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"shanghai" gorm:"index;column:updated_at"`
		DeletedMark int       `json:"deleted_mark"`
		UserId      string    `json:"userId"`
		Name        string    `json:"name"`
	}
	
	dial := sqlite.Open("gorm.db")
	db, _ := gorm.Open(dial, &gorm.Config{})
	
	// 配置待同步的接口
	api := sdk.APIConfig{
		APIPath:   "/api/v1/sync/fakewithts",
		PageSize:  2000,
		BatchSize: 100,
	}

	fakeRows := []FakeRowsWithTS{}

	api.SetParam("ts", "0")
	sdk.SyncToDB(db, api, &fakeRows)

```


更多用法详见以下示例代码，和示例代码中的相关注释

- [Init & CallAPI](example/example.go)
- [SyncToCSV](example/example_csv.go)
- [SyncToModel](example/example_model.go)
- [SyncToDB](example/example_db.go)


## 性能

性能与 ORM 的实现方式（特别是对 upsert 的实现方式），数据库的实现方式，以及网络环境有关，不一定适用于所有情况。

当同步到数据库时，SDK 会采用分批读取/写入的方式，以减少内存的占用。

当同步到模型时，则会将所有数据写入到一个数组中，可能会占用较大的内存。

以下是测试环境

### 同步程序运行环境
 - 4 cpu
 - 8G 内存
 - anolis8 arm64
 - ESSD 云盘 PL0
 - golang 1.21

### 数据库运行环境
mysql/postgresql/sqlserver all in one
 - 2 cpu
 - 8G 内存
 - anolis8 amd64
 - ESSD 云盘 PL0

  
### 测试接口信息
 - /api/v1/sync/fake
 - 使用 pageSize=2000 仅限同步
 - 接口请求耗时约 0.1 - 0.2 秒
 - 接口数据示例

```json
{
	"errCode": 0,
	"errMsg": "success",
	"requestId": "73a60094-c0f1-4daf-bc58-4626fbef7a2b",
	"data": {
		"pageSize": 2000,
		"pageNum": 1,
		"totalNum": 10000,
		"rows": [{
			"id": 1,
			"colString1": "Oxqmn5MWCt",
			"colString2": "mzavQncWeNlOlFgUW7HC",
			"colString3": "mvy6K1HU7rdCicPbvvA3rNZcDWPhvV",
			"colString4": "XGsK5NVQHOu4JrmHZ9ZL1iLf0UYpdIvNIzswULzb",
			"colInt1": 3931594532918648027,
			"colInt2": 337586114254574578,
			"colInt3": 2291922259603323213,
			"colInt4": 3000562485500051124,
			"colFloat1": 0.46541339000557547,
			"colFloat2": 0.6307996439929248,
			"colFloat3": 0.9278393850101392,
			"colFloat4": 0.7286866920659677,
			"colSqlTime1": "2023-10-20 22:02:07",
			"colSqlTime2": "2023-10-20 22:02:07",
			"colSqlTime3": "2023-10-20 22:02:07",
			"colSqlTime4": "2023-10-20 22:02:07"
		}]
	}
}
```

### 10000 数据量
#### 同步到模型

- 耗时：0.68 秒
- 内存分配：30M

#### 同步到数据库

含模型同步时间

|数据库|服务版本|全量写入耗时|全量更新耗时|内存分配|
|--|--|--|--|--|
|MySQL|8.0.34|1.48秒|1.20秒|1.6M|
|PostgreSQL|10.23|0.88秒|0.99秒|2.5M|
|SQLServer|16.0|6.28秒|5.67秒|6M|
|SQLite|/|1.03秒|0.98秒|2.1M|

### 100000 数据量
#### 同步到模型

- 耗时：6.4 秒
- 内存分配：371M

#### 同步到数据库

含模型同步时间

|数据库|服务版本|全量写入耗时|全量更新耗时|内存分配|
|--|--|--|--|--|
|MySQL|8.0.34|11.98秒|12.96秒|2.9M|
|PostgreSQL|10.23|8.34秒|8.33秒|2.5M|
|SQLServer|16.0|55.73秒|54.45秒|4.5M|
|SQLite|/|9.53秒|9.02秒|2.5M|

### 1000000 数据量
#### 同步到模型

- 耗时：63.7 秒
- 内存分配：2773M

#### 同步到数据库

含模型同步时间

|数据库|服务版本|全量写入耗时|全量更新耗时|内存分配|
|--|--|--|--|--|
|MySQL|8.0.34|119.72秒|133.96秒|3.6M|
|PostgreSQL|10.23|86.64秒|90.64秒|3.6M|
|SQLServer|16.0|543.73秒|545.45秒|2.8M|
|SQLite|/|96.6秒|92.62秒|2.2M|
