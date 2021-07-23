package main

import "fmt"

// Data data of a row
type Data map[string]string

// DataMap the data map to store all data
type DataMap map[string]*Data

// Reader gets data from target and return DataMap
type Reader interface {
	Read() *DataMap
}

type Comparer interface {
	Compare(x, y *DataMap)
}

func main() {
	var keys []int
	keys = append(keys, 1)
	r := &ExcelReader{
		filepath:   "t.xlsx",
		keyIndexes: keys,
	}
	dataMap := r.Read()
	for s, data := range *dataMap {
		fmt.Println(s, data)
	}
}
