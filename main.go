package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func main() {
	fares := readDataAndCalculateFares(inputCsvPath)
	file, err := os.Create(outputCsvPath)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for f := range fares {
		var data = []string{strconv.Itoa(f.rideID), strconv.FormatFloat(f.fare, 'f', -1, 64)}
		writer.Write(data)
	}
}

func readDataAndCalculateFares(file string) <-chan Fare {
	out := make(chan Fare, channelBufferSize)
	go func() {
		content, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer content.Close()
		b := bufio.NewReaderSize(content, readerBufferSize)
		r := csv.NewReader(b)
		records := emitStructuredRecords(r)
		groupedRec := groupUniqueRides(records)
		for re := range groupedRec {
			filtered := filterInvalidPoints(re)
			fare := <- estimateFare(filtered)
			out <- fare
		}
		close(out)
	}()
	return out
}

func checkError(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}
