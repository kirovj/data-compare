package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tealeg/xlsx"
)

// List is a list of interface{}
type List []interface{}

// DataMap the data map to store all data
type DataMap map[string]*List

// Result compare result
type Result []*List

// Reader gets data from target and return DataMap
type Reader interface {
	Read() (*DataMap, *List)
}

type equal func(interface{}, interface{}) bool

func basicEqual(x, y interface{}) bool {
	return x == y
}

func Compare(xMap, yMap *DataMap, e equal) (*Result, *Result, *Result) {

	var (
		result Result
		xOnly  Result
		yOnly  Result
	)

	for key, xData := range *xMap {
		yData := (*yMap)[key]

		var r List
		r = append(r, key)

		// only x has
		if yData == nil {
			*xData = append(*xData, key)
			xOnly = append(xOnly, xData)
			continue
		}
		delete(*yMap, key)

		for name, x := range *xData {
			y := (*yData)[name]
			if e(x, y) {
				r = append(r, x, y, "T")
			} else {
				r = append(r, x, y, "F")
			}
		}
		result = append(result, &r)
	}

	// only y has
	for key, yData := range *yMap {
		*yData = append(*yData, key)
		yOnly = append(yOnly, yData)
	}
	return &result, &xOnly, &yOnly
}

// writeExcel write result to xlsx
func writeExcel(result, xOnly, yOnly *Result, cols *List) {

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Difference")

	header := sheet.AddRow()
	for _, col := range *cols {
		cell := header.AddCell()
		cell.Value = col.(string)
	}

	for _, rowData := range *result {
		row := sheet.AddRow()
		for _, val := range *rowData {
			row.AddCell().Value = val.(string)
		}
	}
	xOnlySheet, _ := file.AddSheet("only_X_has")
	for _, rowData := range *xOnly {
		row := xOnlySheet.AddRow()
		for _, val := range *rowData {
			row.AddCell().Value = val.(string)
		}
	}
	yOnlySheet, _ := file.AddSheet("only_Y_has")
	for _, rowData := range *yOnly {
		row := yOnlySheet.AddRow()
		for _, val := range *rowData {
			row.AddCell().Value = val.(string)
		}
	}

	_ = file.Save("result.xlsx")
}

func main() {
	var keys []int

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			k, _ := strconv.Atoi(arg)
			keys = append(keys, k)
		}
	} else {
		keys = append(keys, 1)
	}

	reader := &ExcelReader{
		Filepath:   "x.xlsx",
		KeyIndexes: keys,
	}
	xDataMap, xCols := reader.Read()
	reader.Filepath = "y.xlsx"
	yDataMap, yCols := reader.Read()

	if len(*xCols) != len(*yCols) {
		fmt.Println("column num of two xlsx is different, please make sure they are equal.")
		select {}
	}
	result, xOnly, yOnly := Compare(xDataMap, yDataMap, basicEqual)
	writeExcel(result, xOnly, yOnly, xCols)
}
