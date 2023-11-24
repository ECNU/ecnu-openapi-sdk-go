package sdk

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	xlsx "github.com/tealeg/xlsx/v3"
	"os"
	"reflect"
	"sort"
	"time"

	jsontime "github.com/liamylian/jsontime/v2/v2"
)

var jsonTime = jsontime.ConfigWithCustomTimeFormat

func init() {
	timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
	jsontime.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
	jsontime.AddLocaleAlias("shanghai", timeZoneShanghai)
}

// ParseRowsToXSLX
func parseRowsToXLSX(rows []interface{}, filename string) error {
	// 首先，检查 rows 是否为空
	if len(rows) == 0 {
		return errors.New("rows is empty")
	}

	// 然后，检查 filename 是否为空
	if filename == "" {
		return errors.New("filename is empty")
	}

	// 创建一个新的xlsx文件
	xlsxFile := xlsx.NewFile()

	// 创建一个新的sheet
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	// 创建一个新的行
	xlsxRow := sheet.AddRow()

	// 遍历 rows 中的每个元素，将其转换为 map[string]interface{}
	for i, row := range rows {
		// 将 row 转换为 json 字节
		jsonBytes, err := json.Marshal(row)
		if err != nil {
			return err
		}

		// 创建一个新的 map[string]interface{}
		m := make(map[string]interface{})

		// 解码 json 字节到新 map 中
		err = json.Unmarshal(jsonBytes, &m)
		if err != nil {
			return err
		}

		// 对 map 中的 key 进行排序
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// 对于第一行，写入 header
		if i == 0 {
			for _, key := range keys {
				cell := xlsxRow.AddCell()
				cell.SetString(key)
			}
		}

		// 创建一个新的行
		xlsxRow = sheet.AddRow()

		// 将 values 写入xlsx
		for _, key := range keys {
			cell := xlsxRow.AddCell()
			cell.SetString(fmt.Sprintf("%v", m[key]))
		}
	}

	// 保存xlsx文件
	err = xlsxFile.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

// ParseRowsToCSV
func ParseRowsToCSV(rows []interface{}, filename string) error {
	// 首先，检查 rows 是否为空
	if len(rows) == 0 {
		return errors.New("rows is empty")
	}

	// 然后，检查 filename 是否为空
	if filename == "" {
		return errors.New("filename is empty")
	}

	// 接下来，将 rows 转换为 csv
	// 创建一个新的文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个 csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 遍历 rows 中的每个元素，将其转换为 map[string]interface{}
	for i, row := range rows {
		// 将 row 转换为 json 字节
		jsonBytes, err := json.Marshal(row)
		if err != nil {
			return err
		}
		// 创建一个新的 map[string]interface{}
		m := make(map[string]interface{})
		// 解码 json 字节到新 map 中
		err = json.Unmarshal(jsonBytes, &m)
		if err != nil {
			return err
		}

		var keys, values []string
		// 遍历 map 中的 key，将其放入 keys 中，然后对 keys 排序。
		for k := range m {
			keys = append(keys, k)
			//对keys进行排序
			sort.Strings(keys)
		}
		// 遍历 keys，将 map 中的 value 放入 values 中
		for _, k := range keys {
			values = append(values, fmt.Sprintf("%v", m[k]))
		}

		//对于第一行，先写入 header
		if i == 0 {
			// 将 keys 写入 csv
			err = writer.Write(keys)
			if err != nil {
				return err
			}

		}
		// 将 values 写入 csv
		err = writer.Write(values)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalRows 将一个 []interface{} 的数据映射到一个 struct 数组
func UnmarshalRows(src []interface{}, dst interface{}) error {
	// 首先，检查dst是否是一个指向切片的指针
	dstVal := reflect.ValueOf(dst)

	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return errors.New("dst must be a non-nil pointer to a slice")
	}
	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Slice {
		fmt.Println(dstElem.Kind())
		return errors.New("dst must be a pointer to a slice")
	}

	// 然后，遍历src中的每个元素，将其转换为json字节，并解码到dstElem的元素类型的值中
	for _, s := range src {
		// 将s转换为json字节
		jsonBytes, err := json.Marshal(s)
		if err != nil {
			return err
		}
		// 创建一个新的dstElem的元素类型的值
		newVal := reflect.New(dstElem.Type().Elem())
		// 解码json字节到新值中
		err = jsonTime.Unmarshal(jsonBytes, newVal.Interface())
		if err != nil {
			return err
		}
		// 将新值追加到dstElem中
		dstElem.Set(reflect.Append(dstElem, newVal.Elem()))
	}
	return nil
}

func newStructSlice(input interface{}) (output interface{}, err error) {
	// Dereference input if it is a pointer
	inputValue := reflect.Indirect(reflect.ValueOf(input))

	// Check if input is a slice
	inputType := inputValue.Type()
	if inputType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice")
	}

	// Check if input's element type is a struct
	elemType := inputType.Elem()
	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input's element type is not a struct")
	}

	outputValue := reflect.New(inputValue.Type())
	// Return the output as an interface array
	return outputValue.Interface(), nil
}

func newStruct(input interface{}) (output interface{}, err error) {
	v := reflect.Indirect(reflect.ValueOf(input))
	vType := v.Type()
	if vType.Kind() != reflect.Struct {
		err = fmt.Errorf("input must be struct %v", vType.Kind())
		return
	}

	if !v.IsValid() {
		err = fmt.Errorf("input is not valid %v", v)
		return
	}
	output = reflect.New(vType).Interface()
	return
}
