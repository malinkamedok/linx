package main

import (
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	mmap2 "github.com/edsrzf/mmap-go"
	"log"
	"os"
	"sync"
)

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

	f, _ := os.OpenFile(filename, os.O_RDWR, 0644)
	defer f.Close()

	mmap, _ := mmap2.Map(f, mmap2.RDWR, 0)
	defer mmap.Unmap()

	if filename[len(filename)-5:] == ".json" {
		maxPrice, maxRating, maxPriceName, maxRatingName := jsonParser(mmap)

		fmt.Println("Самый дорогой продукт: ", maxPriceName, " с ценой: ", maxPrice)
		fmt.Println("Продукт с самым высоким рейтингом: ", maxRatingName, " с рейтингом: ", maxRating)
	}

}
