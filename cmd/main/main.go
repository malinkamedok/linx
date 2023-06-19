package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"log"
	"sync"
)

func main() {
	data := []byte(`[
    {"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 10},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5},
	{"product": "Печенье", "price": 34, "rating": 3},
    {"product": "Сахар", "price": 45, "rating": 2},
    {"product": "Варенье", "price": 200, "rating": 5}
]`)

	var maxPrice, maxRating int64
	var maxPriceName, maxRatingName string

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

	fmt.Println("Самый дорогой продукт: ", maxPriceName, " с ценой: ", maxPrice)
	fmt.Println("Продукт с самым высоким рейтингом: ", maxRatingName, " с рейтингом: ", maxRating)

}
