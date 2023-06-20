package main

import (
	csv2 "encoding/csv"
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	mmap2 "github.com/edsrzf/mmap-go"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

func convertStringToInt(inputString <-chan string, outputInt chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(outputInt)

	for strVal := range inputString {
		val, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			fmt.Println("error in parsing int in csv record: ", err)
		}
		outputInt <- val
	}
}

func csvParser(data *os.File) (maxPrice, maxRating int64, maxPriceName, maxRatingName string) {
	reader := csv2.NewReader(data)
	reader.FieldsPerRecord = 3

	wg := new(sync.WaitGroup)
	wg.Add(2)

	//stopChan := make(chan struct{})

	sendStringPrice := make(chan string)

	sendStringRating := make(chan string)

	collectIntPrice := make(chan int64)

	collectIntRating := make(chan int64)

	go convertStringToInt(sendStringPrice, collectIntPrice, wg)
	go convertStringToInt(sendStringRating, collectIntRating, wg)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error in reading record in csv file: ", err)
		}

		if record[1] != "Price" && record[2] != "Rating" {
			sendStringPrice <- record[1]
			sendStringRating <- record[2]
			valPrice := <-collectIntPrice
			valRating := <-collectIntRating

			if valPrice > maxPrice {
				maxPrice = valPrice
				maxPriceName = record[0]
			}

			if valRating > maxRating {
				maxRating = valRating
				maxRatingName = record[0]
			}
		}
	}

	close(sendStringRating)
	close(sendStringPrice)

	wg.Wait()
	return maxPrice, maxRating, maxPriceName, maxRatingName
}

func jsonParser(data []byte) (maxPrice, maxRating int64, maxPriceName, maxRatingName string) {

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			val, err := jsonparser.GetInt(value, "price")
			if err != nil {
				log.Println("error in parsing price")
				return
			}
			if val > maxPrice {
				maxPrice = val
				maxPriceName, err = jsonparser.GetString(value, "product")
				if err != nil {
					log.Println("error in parsing product name with the highest price")
					return
				}
			}
			//fmt.Println(val)
		})
		if err != nil {
			log.Println("error in parsing json")
			return
		}
	}()
	go func() {
		defer wg.Done()
		_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			val, err := jsonparser.GetInt(value, "rating")
			if err != nil {
				log.Println("error in parsing rating")
				return
			}
			if val > maxRating {
				maxRating = val
				maxRatingName, err = jsonparser.GetString(value, "product")
				if err != nil {
					log.Println("error in parsing product name with the highest price")
					return
				}
			}
			//fmt.Println(val)
		})
		if err != nil {
			log.Println("error in parsing json")
			return
		}
	}()

	wg.Wait()
	return maxPrice, maxRating, maxPriceName, maxRatingName
}

func main() {

	//Чтение флага
	var filename string
	flag.StringVar(&filename, "filename", "", "enter your file name")

	flag.Parse()

	fmt.Println("parsed flag: ", filename)

	if filename == "" {
		fmt.Println("Пожалуйста, введите название файла с ключом -filename")
		return
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("error in opening file")
		fmt.Println("Возможно, введен некорректный формат файла. Пожалуйста, введите .json или .csv файл")
		return
	}
	defer f.Close()

	if filename[len(filename)-5:] == ".json" {
		mmapData, err := mmap2.Map(f, mmap2.RDONLY, 0)
		if err != nil {
			fmt.Println("error in mmapping the file: ", err)
			return
		}
		defer mmapData.Unmap()

		maxPrice, maxRating, maxPriceName, maxRatingName := jsonParser(mmapData)

		fmt.Println("Самый дорогой продукт: ", maxPriceName, " с ценой: ", maxPrice)
		fmt.Println("Продукт с самым высоким рейтингом: ", maxRatingName, " с рейтингом: ", maxRating)
	} else if filename[len(filename)-4:] == ".csv" {
		maxPrice, maxRating, maxPriceName, maxRatingName := csvParser(f)

		fmt.Println("Самый дорогой продукт: ", maxPriceName, " с ценой: ", maxPrice)
		fmt.Println("Продукт с самым высоким рейтингом: ", maxRatingName, " с рейтингом: ", maxRating)
	}
}
