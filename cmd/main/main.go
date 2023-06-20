package main

import (
	csv2 "encoding/csv"
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	mmap2 "github.com/edsrzf/mmap-go"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"sync"
)

func convertStringToInt(inputString <-chan string, outputInt chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(outputInt)

	log.Debug("goroutine parsing int started")

	for strVal := range inputString {
		val, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			fmt.Println("error in parsing int in csv record: ", err)
		}
		outputInt <- val
	}

	log.Debug("receiving channel is empty. goroutine exits. send int channel is closing now")
}

func csvParser(data *os.File) (maxPrice, maxRating int64, maxPriceName, maxRatingName string) {
	reader := csv2.NewReader(data)
	reader.FieldsPerRecord = 3

	log.Debug("csv reader init")

	wg := new(sync.WaitGroup)
	wg.Add(2)

	log.Debug("wait group init")

	//stopChan := make(chan struct{})

	sendStringPrice := make(chan string)

	sendStringRating := make(chan string)

	collectIntPrice := make(chan int64)

	collectIntRating := make(chan int64)

	log.Debug("all channels are made")

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
				log.Debug("max price value found")
				maxPrice = valPrice
				maxPriceName = record[0]
			}

			if valRating > maxRating {
				log.Debug("max rating value found")
				maxRating = valRating
				maxRatingName = record[0]
			}
		}
	}

	close(sendStringRating)
	close(sendStringPrice)

	log.Debug("send string channels are closed")

	wg.Wait()

	return maxPrice, maxRating, maxPriceName, maxRatingName
}

func jsonParser(data []byte) (maxPrice, maxRating int64, maxPriceName, maxRatingName string) {

	wg := new(sync.WaitGroup)
	wg.Add(2)

	log.Debug("wait group init")

	go func() {
		defer wg.Done()
		log.Debug("goroutine parsing price started")
		_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			val, err := jsonparser.GetInt(value, "price")
			if err != nil {
				log.Println("error in parsing price")
				return
			}
			if val > maxPrice {
				log.Debug("max price value found")
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
		log.Debug("end of parsing price")
	}()
	go func() {
		defer wg.Done()
		log.Debug("goroutine parsing rating started")
		_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			val, err := jsonparser.GetInt(value, "rating")
			if err != nil {
				log.Println("error in parsing rating")
				return
			}
			if val > maxRating {
				log.Debug("max rating value found")
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
		log.Debug("end of parsing rating")
	}()

	wg.Wait()
	log.Debug("all goroutines are done")
	return maxPrice, maxRating, maxPriceName, maxRatingName
}

func main() {

	log.SetOutput(os.Stdout)

	//Чтение флага
	var filename string
	flag.StringVar(&filename, "filename", "", "enter your file name")

	var logVerbosity bool
	flag.BoolVar(&logVerbosity, "v", false, "set flag true to see debug logs")

	flag.Parse()

	if logVerbosity {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("flags scanned: ", filename, " ", logVerbosity)
	//fmt.Println("parsed flag: ", filename)

	if filename == "" {
		fmt.Println("Пожалуйста, введите название файла с ключом --filename")
		return
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Errorf("error in opening file %v", err)
		fmt.Println("Возможно, введен некорректный формат файла. Пожалуйста, введите .json или .csv файл")
		return
	}
	defer f.Close()

	log.Debug("file is opened")

	if filename[len(filename)-5:] == ".json" {

		log.Debug("starting parsing .json file")

		mmapData, err := mmap2.Map(f, mmap2.RDONLY, 0)
		if err != nil {
			fmt.Println("error in mmapping the file: ", err)
			return
		}
		defer mmapData.Unmap()

		log.Debug("file is mmaped")

		maxPrice, maxRating, maxPriceName, maxRatingName := jsonParser(mmapData)

		log.Debug("end of parsing .json")

		fmt.Println("Самый дорогой продукт: ", maxPriceName, " с ценой: ", maxPrice)
		fmt.Println("Продукт с самым высоким рейтингом: ", maxRatingName, " с рейтингом: ", maxRating)
	} else if filename[len(filename)-4:] == ".csv" {
		log.Debug("starting parsing .csv file")

		maxPrice, maxRating, maxPriceName, maxRatingName := csvParser(f)

		log.Debug("end of parsing .csv")

		fmt.Println("Самый дорогой продукт: ", maxPriceName, " с ценой: ", maxPrice)
		fmt.Println("Продукт с самым высоким рейтингом: ", maxRatingName, " с рейтингом: ", maxRating)
	}
}
