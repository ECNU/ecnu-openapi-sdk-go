package main

import (
	"fmt"

	"github.com/ecnu/ecnu-openapi-sdk-go/sdk"
)

func exampleSyncToModel() {
	// 全量同步有效的部分
	api := sdk.APIConfig{
		APIPath:  "/api/v1/sync/fakewithts",
		PageSize: 2000,
	}
	api.SetParam("ts", "0")
	fakeRows := []FakeRowsWithTS{}
	if err := sdk.SyncToModel(api, &fakeRows); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Model：全量同步 %d 条数据\n", len(fakeRows))

	//可以根据自己的需要处理数据
	/*
		for _, row := range fakeRows{
			do something
		}
	*/

	// 增量同步，返回 2023-01-03 00:00:00 之后的全部数据（含软删除部分）
	ts := 1672675200

	api.SetParam("ts", fmt.Sprintf("%d", ts))
	api.SetParam("full", "1")

	fakeRows = []FakeRowsWithTS{}
	if err := sdk.SyncToModel(api, &fakeRows); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Model:增量同步 %d 条数据\n", len(fakeRows))

}
