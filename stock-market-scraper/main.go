package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

func main() {
	tickers := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := []Stock{}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("div.region--intraday", func(e *colly.HTMLElement) {
		stock := Stock{}

		// Extract the company name from the URL or use a static value if preferred
		stock.company = e.Request.URL.Query().Get("stock")
		if stock.company == "" {
			stock.company = strings.ToUpper(strings.Split(e.Request.URL.Path, "/")[3])
		}

		stock.price = e.ChildText("bg-quote.value")
		fmt.Println("Price:", stock.price)
		stock.change = e.ChildText("span.change--percent--q")
		fmt.Println("Change:", stock.change)

		stocks = append(stocks, stock)
	})

	for _, t := range tickers {
		c.Visit("https://www.marketwatch.com/investing/stock/" + t)
	}

	c.Wait() // Wait for all requests to complete

	fmt.Println("Extracted stocks:", stocks)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Println("Error creating 'stocks.csv' file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	headers := []string{
		"company",
		"price",
		"change",
	}
	writer.Write(headers)

	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}

	writer.Flush() // Ensure all buffered data is written to the file
}
