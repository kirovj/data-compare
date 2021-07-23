package main

import (
	"github.com/tealeg/xlsx"
)

type ExcelReader struct {
	filepath   string // excel file path
	keyIndexes []int  // keyIndexes is the key column index to find a data
}

func (e *ExcelReader) Read() *DataMap {
	var (
		err   error
		excel *xlsx.File
		cols  []string
	)
	if excel, err = xlsx.OpenFile(e.filepath); err != nil {
		return nil
	}

	sheet := excel.Sheets[0]
	head := sheet.Row(0)
	for _, cell := range head.Cells {
		cols = append(cols, cell.String())
	}

	var dataMap = make(DataMap)

	for _, row := range sheet.Rows[1:] {
		var rowData = make(Data)
		key := ""
		for i, cell := range row.Cells {

			// 判断是否为主键列
			for _, index := range e.keyIndexes {
				if index == i+1 {
					key += cell.String()
				}
				continue
			}
			rowData[cols[i]] = cell.String()
		}
		dataMap[key] = &rowData
	}

	return &dataMap
}
