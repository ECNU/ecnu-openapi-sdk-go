package main

import (
	"fmt"

	"github.com/ecnu/ecnu-openapi-sdk-go/sdk"
)

func exampleSyncToXLSX() {
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
		APIPath:  "/api/v1/sync/fakewithts",
		PageSize: 2000,
	}

	api.SetParam("ts", "0")
	// 同步到 xlsx
	// xlsx 模式下，所有字段都会转为 string
	xlsxFile := "test.xlsx"
	mode := "xlsx"
	rows, err := sdk.SyncToFile(mode, xlsxFile, api)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("XLSX：组织机构同步 %d 条数据\n", rows)
}
