package main

import (
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
	"os"
)

func main() {

	filepath := os.Args[1]
	jsonstr := os.Args[2]

	fmt.Println(jsonstr)

	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := f.GetRows("MHKN_D2C")
	if err != nil {
		fmt.Println(err)
		return
	}

	records := [][]string{}

	for _, row := range rows {
		line := []string{}
		for _, colCell := range row {
			line = append(line, colCell)
		}
		records = append(records, line)
	}

	output, err := os.Create("result.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer output.Close()

	w := csv.NewWriter(output)
	defer w.Flush()

	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
