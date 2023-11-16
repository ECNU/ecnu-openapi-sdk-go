package main

import "time"

/*
		根据实际的接口字段，定义模型
		模型的定义约定，详见 gorm 相关文档：https://gorm.io/zh_CN/docs/models.html

		建议约定好主键，避免数据重复
		如果不设定主键，则每次同步前需要清空表，否则会产生重复数据

	    由于 gorm 默认会使用 CreatedAt、UpdatedAt 追踪创建/更新时间，DeletedAt 追踪删除时间。
	    如果接口返回的内容涉及相关含义的字段，则在定义模型时，需要避开这几个字段名，以避免被 gorm 自动更新了相关字段导致增量同步失败（特别是 UpdatedAt)。
	    例如将接口返回的 updatedAt，在模型中使用 UpdateTime 作为字段，而非 UpdatedAt。
		然而到数据库字段间的映射，我们依然可以指定他的名称，例如：gorm:"index;column:updated_at"
		详见：https://gorm.io/zh_CN/docs/models.html

		gorm 中默认整形的主键会开启 AutoIncrement，这可能会导致数据不一致而产生同步异常（特别是sqlserver等常见）
		因此如果字段的主键是整形，请务必显式的关闭自增，例如：gorm:"primarykey;autoIncrement:false"

		sql datetime 格式的字符串，直接解析到 time.Time，需要使用 time_format 和 time_location 两个 tag
		详见，github.com/liamylian/jsontime/v2/v2
*/

type FakeRowsWithTS struct {
	Id          int       `json:"id" gorm:"primarykey;autoIncrement:false"`
	CreateTime  time.Time `json:"created_at" time_format:"sql_datetime" time_location:"shanghai"`
	UpdateTime  time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"shanghai" gorm:"index;column:updated_at"`
	DeletedMark int       `json:"deleted_mark"`
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
}

type FakeRows struct {
	Id          int       `json:"id" gorm:"primarykey;autoIncrement:false"`
	ColString1  string    `json:"colString1"`
	ColString2  string    `json:"colString2"`
	ColString3  string    `json:"colString3"`
	ColString4  string    `json:"colString4"`
	ColInt1     int64     `json:"colInt1"`
	ColInt2     int64     `json:"colInt2"`
	ColInt3     int64     `json:"colInt3"`
	ColInt4     int64     `json:"colInt4"`
	ColFloat1   float64   `json:"colFloat1"`
	ColFloat2   float64   `json:"colFloat2"`
	ColFloat3   float64   `json:"colFloat3"`
	ColFloat4   float64   `json:"colFloat4"`
	ColSqlTime1 time.Time `json:"colSqlTime1" time_format:"sql_datetime" time_location:"shanghai"`
	ColSqlTime2 time.Time `json:"colSqlTime2" time_format:"sql_datetime" time_location:"shanghai"`
	ColSqlTime3 time.Time `json:"colSqlTime3" time_format:"sql_datetime" time_location:"shanghai"`
	ColSqlTime4 time.Time `json:"colSqlTime4" time_format:"sql_datetime" time_location:"shanghai"`
}
