package main

import (
	"encoding/csv"
	"log"
	"os"
)

func ExportCsv(data [][]string, csvFileName string) error {
	file, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.WriteAll(data)

	if err = writer.Error(); err != nil {
		log.Fatalln("csv writing error:", err)
	}
	return nil
}
