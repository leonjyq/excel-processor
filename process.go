package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Channel struct {
	Channel    string `json:"channel"`
	Percentage int    `json:"percentage"`
}

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
	const inForm = "Jan-06"
	const outForm = "01/02/2006"
	channels := []Channel{}
	records := [][]string{}
	header := []string{"Transaction Date", "Value in document currency", "Value in Company Currency", "Units Sold",
		"Item number affiliate", "Customer number local", "Posting Date"}
	records = append(records, header)

	for _, sheet := range sheets {
		fmt.Println(sheet)
		sheetType := strings.Split(sheet, "_")[1]
		if strings.EqualFold(sheetType, "distributor") {
			json.Unmarshal([]byte(distributorRule), &channels)
		} else if strings.EqualFold(sheetType, "d2c") {
			json.Unmarshal([]byte(d2cRule), &channels)
		}
		cols, err := f.GetCols(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}
		date := cols[1][1]
		fmt.Println(date)
		t, _ := time.Parse(inForm, date)
		m := int(t.Month())
		restMonth := 12 - m + 1
		for i, col := range cols[1 : restMonth+1] {
			// kind := col[0]
			date := col[1]
			t, _ := time.Parse(inForm, date)
			last := t.AddDate(0, 1, -1)
			transactionDate := last.Format(outForm)
			for j, colValue := range col[2:] {
				for _, channel := range channels {
					line := []string{}
					line = append(line, transactionDate)
					fmt.Print(i + 1 + restMonth)
					fmt.Print("--")
					fmt.Println(j + 2)
					fmt.Print(cols[i+1+restMonth][j+2])
					fmt.Print("--")
					fmt.Println(colValue)
					revenue, err := strconv.Atoi(cols[i+1+restMonth][j+2])

					if err != nil {
						fmt.Println(err)
						return
					}

					unitSold, err := strconv.Atoi(colValue)
					if err != nil {
						fmt.Println(err)
						return
					}
					line = append(line, strconv.Itoa(revenue*channel.Percentage/100))
					line = append(line, strconv.Itoa(revenue*channel.Percentage/100))
					line = append(line, strconv.Itoa(unitSold*channel.Percentage/100))
					line = append(line, cols[0][j])
					line = append(line, channel.Channel)
					line = append(line, transactionDate)
					records = append(records, line)
				}
			}
		}
	}

	return

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
