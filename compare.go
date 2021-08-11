package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tealeg/xlsx"
)

// List is a list of string
type List []string

// DataMap the data map to store all data
type DataMap map[string]*List

// Result compare result
type Result []*List

// Reader gets data from target and return DataMap
type Reader interface {
	Read() (*DataMap, *List)
}

type Comparer interface {
	Compare(x, y *DataMap) (*Result, *Result, *Result)
	IsEqual(x, y string) bool
}

type CommonCompare struct {
	Round uint8 // float round num
}

func (c *CommonCompare) IsEqual(x, y string) bool {
	return x == y
}

func (c *CommonCompare) Compare(xMap, yMap *DataMap) (*Result, *Result, *Result) {

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
			if c.IsEqual(x, y) {
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
	sheet, _ := file.AddSheet("Sheet1")

	header := sheet.AddRow()
	for _, col := range *cols {
		cell := header.AddCell()
		cell.Value = col
	}

	for _, rowData := range *result {
		row := sheet.AddRow()
		for _, val := range *rowData {
			row.AddCell().Value = val
		}
	}
	xOnlySheet, _ := file.AddSheet("xOnly")
	for _, rowData := range *xOnly {
		row := xOnlySheet.AddRow()
		for _, val := range *rowData {
			row.AddCell().Value = val
		}
	}
	yOnlySheet, _ := file.AddSheet("yOnly")
	for _, rowData := range *yOnly {
		row := yOnlySheet.AddRow()
		for _, val := range *rowData {
			row.AddCell().Value = val
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
		fmt.Println("column num of two xlsx is different, please mark sure they are equal.")
		select {}
	}

	c := &CommonCompare{}
	result, xOnly, yOnly := c.Compare(xDataMap, yDataMap)
	writeExcel(result, xOnly, yOnly, xCols)
}
