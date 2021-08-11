package main

import (
	"strings"

	"github.com/tealeg/xlsx"
)

type ExcelReader struct {
	Filepath   string // excel file path
	KeyIndexes []int  // keyIndexes is the key column index to find a data
}

func getSheet(file string) (*xlsx.Sheet, error) {
	excel, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, err
	}
	return excel.Sheets[0], nil
}

func (e *ExcelReader) isKey(i int) bool {
	for _, index := range e.KeyIndexes {
		if index == i+1 {
			return true
		}
	}
	return false
}

func (e *ExcelReader) Read() (*DataMap, *List) {
	var (
		err   error
		sheet *xlsx.Sheet
		cols  List
	)
	if sheet, err = getSheet(e.Filepath); err != nil {
		return nil, nil
	}

	// get headers
	head := sheet.Row(0)
	cols = append(cols, "")
	for i, cell := range head.Cells {
		if e.isKey(i) {
			continue
		}
		name := cell.String()
		cols = append(append(append(cols, name+"_x"), name+"_y"), "result")
	}

	var dataMap = make(DataMap)
	for _, row := range sheet.Rows[1:] {
		var rowData List
		var key []string
		for i, cell := range row.Cells {
			// 判断是否为主键列
			if e.isKey(i) {
				key = append(key, cell.String())
				continue
			}
			rowData = append(rowData, cell.String())
		}
		dataMap[strings.Join(key, "|")] = &rowData
	}
	return &dataMap, &cols
}
