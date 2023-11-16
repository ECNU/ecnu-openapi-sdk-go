package sdk

import (
	"database/sql"
	"net/url"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type APIConfig struct {
	APIPath        string `json:"api_path"`
	PageSize       int    `json:"page_size"`
	BatchSize      int    `json:"batch_size"`
	UpdatedAtField string
	params         url.Values
}

func (api *APIConfig) SetDefault() {
	if api.PageSize == 0 {
		api.PageSize = 2000
	}
	if api.PageSize > MAXPageSIZE {
		api.PageSize = MAXPageSIZE
	}
	if api.BatchSize == 0 {
		api.BatchSize = 100
	}
	if api.UpdatedAtField == "" {
		api.UpdatedAtField = "updated_at"
	}
}

func (api *APIConfig) AddParam(key, value string) {
	if api.params == nil {
		api.params = make(url.Values)
	}
	api.params.Add(key, value)
}
func (api *APIConfig) SetParam(key, value string) {
	if api.params == nil {
		api.params = make(url.Values)
	}
	api.params.Set(key, value)
}

func (api *APIConfig) DelParam(key string) {
	if api.params == nil {
		api.params = make(url.Values)
	}
	api.params.Del(key)
}

func (api *APIConfig) ParamEncode() string {
	return api.params.Encode()
}

func SyncToCSV(csvFileName string, api APIConfig) (int64, error) {
	c := GetOpenAPIClient()
	api.SetDefault()
	apiPath := api.APIPath
	if api.ParamEncode() != "" {
		if strings.Contains(apiPath, "?") {
			apiPath = api.APIPath + "&" + api.ParamEncode()
		} else {
			apiPath = api.APIPath + "?" + api.ParamEncode()
		}
	}
	rows, err := c.GetAllRows(apiPath, api.PageSize)
	if err != nil {
		return 0, err
	}
	if err = ParseRowsToCSV(rows, csvFileName); err != nil {
		return 0, err
	}
	return int64(len(rows)), nil
}

func SyncToModel(api APIConfig, dataModel interface{}) error {
	apiPath := api.APIPath
	if api.ParamEncode() != "" {
		if strings.Contains(apiPath, "?") {
			apiPath = api.APIPath + "&" + api.ParamEncode()
		} else {
			apiPath = api.APIPath + "?" + api.ParamEncode()
		}
	}
	c := GetOpenAPIClient()
	rows, err := c.GetAllRows(apiPath, api.PageSize)
	if err != nil {
		return err
	}
	if err := UnmarshalRows(rows, dataModel); err != nil {
		return err
	}
	return nil
}

func SyncToDB(db *gorm.DB, api APIConfig, dataModel interface{}) (int64, error) {
	api.SetDefault()
	if err := db.AutoMigrate(dataModel); err != nil {
		return 0, err
	}
	apiPath := api.APIPath
	if api.ParamEncode() != "" {
		if strings.Contains(apiPath, "?") {
			apiPath = api.APIPath + "&" + api.ParamEncode()
		} else {
			apiPath = api.APIPath + "?" + api.ParamEncode()
		}
	}
	c := GetOpenAPIClient()
	pageNum := 1
	rowsCount := int64(0)
	for {
		res, err := c.GetRows(apiPath, pageNum, api.PageSize)
		if err != nil {
			return rowsCount, err
		}

		//利用反射创建一个结构相同的临时空间，是个指针
		tmpData, err := newStructSlice(dataModel)
		if err != nil {
			return rowsCount, err
		}

		if err := UnmarshalRows(res.Rows, tmpData); err != nil {
			return rowsCount, err
		}

		//如果空指针后面反射会 panic，容错性处理
		if tmpData == nil {
			break
		} else {
			//数据结构已知，如果不是空指针那一定是数组，所以 Len 方法必然有效
			v := reflect.Indirect(reflect.ValueOf(tmpData))
			if v.Len() == 0 {
				break
			}
			rowsCount = rowsCount + int64(v.Len())
		}

		sqlResult := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(tmpData, api.BatchSize)

		if sqlResult.Error != nil {
			return rowsCount, sqlResult.Error
		}

		pageNum = pageNum + 1
	}
	return rowsCount, nil
}

func GetLastUpdatedTS(db *gorm.DB, api APIConfig, dataModel interface{}) int64 {
	api.SetDefault()
	type result struct {
		TS sql.NullTime `gorm:"column:ts"`
	}
	var res result
	db.Model(&dataModel).Select(api.UpdatedAtField + " as ts").Order(api.UpdatedAtField + " desc").Limit(1).Scan(&res)

	if !res.TS.Valid {
		return 0
	}

	return res.TS.Time.Unix()
}
