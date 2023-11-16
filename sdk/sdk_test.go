package sdk

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func Test_ParseRowsToCSV(t *testing.T) {

	type testStruct struct {
		D string
		C int
		B string
		A float64
	}

	rows := []interface{}{
		testStruct{
			D: "d1",
			C: 1,
			B: "b1",
			A: 1.1,
		},
		testStruct{
			D: "d2",
			C: 2,
			B: "b2",
			A: 2.2,
		},
	}

	filename := "test.csv"

	if err := ParseRowsToCSV(rows, filename); err != nil {
		t.Error(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	/*
		读取 test.csv 文件
		比对第一行的内容是否等于 testStruct 的每个字段名
		比对第二行和第三行的内容是否等于对应的两条数据
	*/
	records, err := reader.ReadAll()
	if err != nil {
		t.Error(err)
	}
	if len(records) != 3 {
		t.Errorf("len(records) != 3")
	}
	// 排序过，应该按 A,B,C,D 排序
	row0 := []string{"A", "B", "C", "D"}
	// 都会被转为字符串
	row1 := []string{"1.1", "b1", "1", "d1"}
	row2 := []string{"2.2", "b2", "2", "d2"}

	for i, record := range records {
		if i == 0 {
			for j, field := range record {
				if field != row0[j] {
					t.Errorf("field != row0[j]")
				}
			}
		} else if i == 1 {
			for j, field := range record {
				if field != row1[j] {
					t.Errorf("field != row1[j]")
				}
			}
		} else if i == 2 {
			for j, field := range record {
				if field != row2[j] {
					t.Errorf("field != row2[j]")
				}
			}
		}
	}
}

func Test_UnmarshalRows(t *testing.T) {
	type testSrcStruct struct {
		A string
		B string
		C string
		D string
	}

	type testDstStruct struct {
		A string
		B string
		C string
		D string
	}

	srcRows := []interface{}{
		testSrcStruct{
			A: "a1",
			B: "b1",
			C: "c1",
			D: "d1",
		},
		testSrcStruct{
			A: "a2",
			B: "b2",
			C: "c2",
			D: "d2",
		},
	}
	dstRows := []testDstStruct{}

	if err := UnmarshalRows(srcRows, &dstRows); err != nil {
		t.Error(err)
	}

	if len(dstRows) != 2 {
		t.Errorf("len(dstRows) != 2")
	}

	for i, row := range dstRows {
		if row.A != srcRows[i].(testSrcStruct).A {
			t.Error("row.A != srcRows[i].A")
		}
		if row.B != srcRows[i].(testSrcStruct).B {
			t.Error("row.B != srcRows[i].B")
		}
		if row.C != srcRows[i].(testSrcStruct).C {
			t.Error("row.C != srcRows[i].C")
		}
		if row.D != srcRows[i].(testSrcStruct).D {
			t.Error("row.D != srcRows[i].D")
		}
	}

}

func Test_newStructSlice(t *testing.T) {
	type testSrcStruct struct {
		A string
		B string
		C string
		D string
	}
	srcRows := []testSrcStruct{
		{
			A: "a1",
			B: "b1",
			C: "c1",
			D: "d1",
		},
		{
			A: "a2",
			B: "b2",
			C: "c2",
			D: "d2",
		},
	}
	dstRows, err := newStructSlice(srcRows)
	if err != nil {
		t.Error(err)
		return
	}
	if fmt.Sprintf("%T", dstRows) != "*[]sdk.testSrcStruct" {
		t.Errorf("dst type is not []sdk.testSrcStruct, current is %T", dstRows)
	}
	jsvalue := `
	[{
		"A":"a1",
		"B":"b1",
		"C":"c1",
		"D":"d1"
	},
	{
		"A":"a2",
		"B":"b2",
		"C":"c2",
		"D":"d2"
	}]
	`
	if err := json.Unmarshal([]byte(jsvalue), &dstRows); err != nil {
		t.Error(err)
		return
	}
	if fmt.Sprintf("%T", dstRows) != "*[]sdk.testSrcStruct" {
		t.Errorf("dst type is not []sdk.testSrcStruct, current is %T", dstRows)
	}

}

func Test_newStruc(t *testing.T) {
	type testSrcStruct struct {
		A string
		B string
		C string
		D string
	}
	src := testSrcStruct{
		A: "a1",
		B: "b1",
		C: "c1",
		D: "d1",
	}

	dst, err := newStruct(src)
	if err != nil {
		t.Error(err)
		return
	}
	if fmt.Sprintf("%T", dst) != "*sdk.testSrcStruct" {
		t.Errorf("dst type is not sdk.testSrcStruct, current is %T", dst)
	}
	jsvalue := `
	{
		"A":"a1",
		"B":"b1",
		"C":"c1",
		"D":"d1"
	}
	`
	if err := json.Unmarshal([]byte(jsvalue), &dst); err != nil {
		t.Error(err)
		return
	}
	if fmt.Sprintf("%T", dst) != "*sdk.testSrcStruct" {
		t.Errorf("dst type is not sdk.testSrcStruct, current is %T", dst)
	}

}
