package main

import (
	"github.com/tealeg/xlsx"
)

// rowData data of a excel row
type rowData map[string]string

// dataMap the data map of one excel
type dataMap map[string]*rowData

type res struct {
	x, y string
	t    uint8
}

func readFile(file string, keyIndexes ...int) *dataMap {
	var (
		err   error
		excel *xlsx.File
		cols  []string
	)
	if excel, err = xlsx.OpenFile(file); err != nil {
		return nil
	}

	sheet := excel.Sheets[0]
	head := sheet.Row(0)
	for _, cell := range head.Cells {
		cols = append(cols, cell.String())
	}

	var dataMap = make(dataMap)

	for _, row := range sheet.Rows[1:] {
		var rowData = make(rowData)
		key := ""
		for i, cell := range row.Cells {

			// 判断是否为主键列
			for _, index := range keyIndexes {
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
