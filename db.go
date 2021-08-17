package main

type DbReader struct {
	Filepath   string // excel file path
	KeyIndexes []int  // keyIndexes is the key column index to find a data
}
