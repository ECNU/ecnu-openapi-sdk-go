package sdk

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DataResult
type DataResult struct {
	TotalNum int           `json:"totalNum"`
	PageSize int           `json:"pageSize"`
	PageNum  int           `json:"pageNum"`
	Rows     []interface{} `json:"rows"`
}

// GetRows
func (c *OAuth2Client) GetRows(apiPath string, pageNum, pageSize int) (DataResult, error) {
	var dataResult DataResult
	var url string
	if strings.Contains(apiPath, "?") {
		url = fmt.Sprintf("%s%s&pageNum=%d&pageSize=%d", c.BaseUrl, apiPath, pageNum, pageSize)
	} else {
		url = fmt.Sprintf("%s%s?pageNum=%d&pageSize=%d", c.BaseUrl, apiPath, pageNum, pageSize)
	}

	data, err := c.HttpGet(url)
	if err != nil {
		return dataResult, err
	}

	if err := json.Unmarshal(data, &dataResult); err != nil {
		return dataResult, err
	}
	return dataResult, err
}

// GetAllRows
func (c *OAuth2Client) GetAllRows(apiPath string, pageSize int) ([]interface{}, error) {
	var rows []interface{}
	pageNum := 1
	for {
		result, err := c.GetRows(apiPath, pageNum, pageSize)
		if err != nil {
			return rows, err
		}
		if len(result.Rows) == 0 {
			break
		}
		pageNum = pageNum + 1
		rows = append(rows, result.Rows...)
	}
	return rows, nil
}
