package csv

import (
	"encoding/csv"
	"os"
)

func LoadFile(path string) (csvRecords [][]string) {
	if csvFile, err := os.Open(path); err != nil {
		panic(err)
	} else {
		defer csvFile.Close()
		records, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			panic(err)
		} else {
			for index, record := range records {
				if index != 0 {
					csvRecords = append(csvRecords, record)
				}
			}
			return csvRecords
		}
	}
}
