package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
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
	const outForm = "02/01/2006"
	channels := []Channel{}
	records := [][]string{}
	header := []string{"Transaction Date", "Value in document currency", "Value in Company Currency", "Units Sold",
		"Item number affiliate", "Customer number local", "Posting Date"}
	records = append(records, header)

	for _, sheet := range sheets {
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
		t, _ := time.Parse(inForm, date)
		m := int(t.Month())
		restMonth := 12 - m + 1
		for i, col := range cols[1 : restMonth+1] {
			i++
			date := col[1]
			t, _ := time.Parse(inForm, date)
			last := t.AddDate(0, 1, -1)
			transactionDate := last.Format(outForm)
			for j, colValue := range col[2:] {
				j = j + 2
				var us, re, fg float64
				for k, channel := range channels {
					line := []string{}
					fgline := []string{}
					line = append(line, transactionDate)
					fgline = append(fgline, transactionDate)
					perc := float64(channel.Percentage) / 100
					freeg, err := strconv.ParseFloat(cols[i+restMonth][j], 64)
					revenue, err := strconv.ParseFloat(cols[i+restMonth+restMonth][j], 64)
					if err != nil {
						// fmt.Print("cell:")
						// fmt.Print(i)
						// fmt.Print("&")
						// fmt.Print(j)
						// fmt.Print("-")
						fmt.Println(err)
					}

					unitSold, err := strconv.ParseFloat(colValue, 64)
					if err != nil {
						// fmt.Print("cell:")
						// fmt.Print(i)
						// fmt.Print("&")
						// fmt.Print(j)
						// fmt.Print("-")
						fmt.Println(err)
					}
					fgline = append(fgline, "0")
					fgline = append(fgline, "0")
					if k != len(channels)-1 {
						re = re + revenue*perc
						us = us + unitSold*perc
						fg = fg + freeg*perc
						// fmt.Print(us)
						// fmt.Print("----")
						// fmt.Println(re)
						line = append(line, strconv.FormatFloat(revenue*perc, 'f', -1, 64))
						line = append(line, strconv.FormatFloat(revenue*perc, 'f', -1, 64))
						line = append(line, strconv.FormatFloat(math.Floor(unitSold*perc+0.5), 'f', -1, 64))
						fgline = append(fgline, strconv.FormatFloat(math.Floor(freeg*perc+0.5), 'f', -1, 64))
						// fmt.Print(colValue)
						// fmt.Print("----")
						// fmt.Print(unitSold)
						// fmt.Print("----")
						// fmt.Print(perc)
						// fmt.Print("----")
						// fmt.Print(unitSold * perc)
						// fmt.Print("----")
						// fmt.Println(math.Floor(unitSold * perc))
					} else {
						line = append(line, strconv.FormatFloat(revenue-re, 'f', -1, 64))
						line = append(line, strconv.FormatFloat(revenue-re, 'f', -1, 64))
						line = append(line, strconv.FormatFloat(math.Floor(unitSold-us+0.5), 'f', -1, 64))
						fgline = append(fgline, strconv.FormatFloat(math.Floor(freeg-fg+0.5), 'f', -1, 64))
					}
					line = append(line, cols[0][j])
					line = append(line, channel.Channel)
					line = append(line, transactionDate)
					fgline = append(fgline, cols[0][j])
					fgline = append(fgline, channel.Channel)
					fgline = append(fgline, transactionDate)
					//fmt.Println(line)
					records = append(records, line)
					records = append(records, fgline)
				}
			}
		}
	}

	output, err := os.Create(outputFile)
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
