package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {

	var inputFile, outputFile, distributorRule, d2cRule string
	flag.StringVar(&inputFile, "in", "", "input file path")
	flag.StringVar(&outputFile, "out", "", "output file path")
	flag.StringVar(&distributorRule, "distributor", "", "rule for Distributor")
	flag.StringVar(&d2cRule, "d2c", "", "rule for D2C")

	flag.Parse()
	fmt.Println(distributorRule)
	fmt.Println(d2cRule)
	fmt.Println(inputFile)
	fmt.Println(outputFile)

	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		fmt.Println(sheet)
		sheetType := strings.Split(sheet, "_")[1]
		if strings.EqualFold(sheetType, "distributor") {
			fmt.Println(sheetType)
		} else if strings.EqualFold(sheetType, "d2c") {
			fmt.Println(sheetType)
		}

		col, err := f.GetCols(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
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
