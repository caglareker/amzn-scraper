package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type product struct {
	Price            string `json:"price"`
	ProductTitle     string `json:"productTitle"`
	ShortDescription string `json:"shortDescription"`
	Reviews          string `json:"reviews"`
}

func main() {
	asins, err := os.Open("asins.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer asins.Close()

	csvReader := csv.NewReader(asins)
	asinCodes, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	urls := make([]string, 0)
	for _, asinCode := range asinCodes {
		url := fmt.Sprintf("https://www.amazon.com/dp/%s", asinCode[0])
		urls = append(urls, url)
	}
	crawl(urls)
}

func crawl(urls []string) {
	allProducts := make([]product, 0)
	replacer := strings.NewReplacer("$", "", ",", "")

	collector := colly.NewCollector(
		colly.AllowedDomains("amazon.com", "www.amazon.com"),
	)

	collector.OnHTML("html", func(e *colly.HTMLElement) {
		model := product{}
		model.Price = replacer.Replace(e.ChildText("#price_inside_buybox"))
		model.ProductTitle = e.ChildText("#productTitle")
		model.ShortDescription = e.ChildText("#featurebullets_feature_div")
		model.Reviews = e.ChildText("#acrCustomerReviewText")

		allProducts = append(allProducts, model)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	for _, url := range urls {
		collector.Visit(url)
	}
	writeJSON(allProducts)
}

func writeJSON(data []product) {
	file, err := json.MarshalIndent(data, "", " ")
	inJSON, _ := _UnescapeUnicodeCharactersInJSON(file)
	if err != nil {
		log.Println("Unable to create JSON file")
		return
	}
	fileName := fmt.Sprintf("products.json")
	_ = ioutil.WriteFile(fileName, inJSON, 0644)
}

func _UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
