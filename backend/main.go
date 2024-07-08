package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/rs/cors"
)

type FinancialData struct {
	Symbol                               string `json:"symbol"`
	Price                                string `json:"price"`
	MarketCap                            string `json:"marketCap"`
	Revenue                              string `json:"revenue"`
	CostOfRevenue                        string `json:"costOfRevenue"`
	GrossProfit                          string `json:"grossProfit"`
	OperatingExpense                     string `json:"operatingExpense"`
	OperatingIncome                      string `json:"operatingIncome"`
	NetNonOperatingInterestIncomeExpense string `json:"netNonOperatingInterestIncomeExpense"`
	OtherIncomeExpense                   string `json:"otherIncomeExpense"`
	PretaxIncome                         string `json:"pretaxIncome"`
	// ... (add other relevant fields from the table)
}

var (
	data = make(map[string]FinancialData)
	mu   sync.RWMutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/compare", handleCompare)

	// Update CORS configuration
	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://financial-comparison-frontend.onrender.com"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	go updateData() // Start the background scraper

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	// w: http.ResponseWriter is for preparing the response to frontend.
	// r: http.Request is for getting the request from frontend.
	w.Header().Set("Content-Type", "application/json")

	// symbols := r.URL.Query()["symbols"]
	// result := make(map[string]FinancialData)

	symbols, ok := r.URL.Query()["symbols[]"]
	if !ok || len(symbols) == 0 {
		http.Error(w, "Missing symbols parameter", http.StatusBadRequest)
		return
	}
	result := make(map[string]FinancialData)

	fmt.Println(symbols)
	fmt.Println(result)

	for _, symbol := range symbols {
		mu.RLock()
		financialData, exists := data[symbol]
		mu.RUnlock()

		if !exists {
			financialData = scrapeFinancialData(symbol)
			mu.Lock()
			data[symbol] = financialData
			mu.Unlock()
		}

		result[symbol] = financialData
	}

	json.NewEncoder(w).Encode(result)
}

func scrapeFinancialData(symbol string) FinancialData {
	// anti-scraping measures
	time.Sleep(2 * time.Second)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	financialData := FinancialData{Symbol: symbol}

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping %s: %v", symbol, err)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("Visited %s, status code: %d", r.Request.URL, r.StatusCode)
	})

	c.OnHTML("fin-streamer[data-field='regularMarketPrice']", func(e *colly.HTMLElement) {
		if financialData.Price == "" {
			financialData.Price = e.Text
			log.Printf("Price for %s: %s", symbol, financialData.Price)
		}
	})

	// c.OnHTML("div.tableBody", func(e *colly.HTMLElement) {
	// 	e.ForEach("div.row.lv-0.svelte-1xjz32c", func(_ int, el *colly.HTMLElement) { // Iterate through the rows
	// 		rowTitle := el.ChildText("div.rowTitle")
	// 		if rowTitle == "Total Revenue" {
	// 			// Find the 3rd column with the class containing revenue data
	// 			revenueText := el.ChildText("div.column.svelte-1xjz32c:nth-child(3)")
	// 			financialData.Revenue = revenueText
	// 			log.Printf("Revenue for %s: %s", symbol, financialData.Revenue)
	// 			return // Stop iterating after finding Total Revenue
	// 		}
	// 	})
	// })

	c.OnHTML("div.tableBody", func(e *colly.HTMLElement) {
		e.ForEach("div.row.lv-0.svelte-1xjz32c", func(_ int, el *colly.HTMLElement) {
			rowTitle := el.ChildText("div.rowTitle")
			valueText := el.ChildText("div.column.svelte-1xjz32c:nth-child(3)")

			switch rowTitle {
			case "Total Revenue":
				financialData.Revenue = valueText
			case "Cost of Revenue":
				financialData.CostOfRevenue = valueText
			case "Gross Profit":
				financialData.GrossProfit = valueText
			case "Operating Expense":
				financialData.OperatingExpense = valueText
			case "Operating Income":
				financialData.OperatingIncome = valueText
			case "Net Non Operating Interest Income Expense":
				financialData.NetNonOperatingInterestIncomeExpense = valueText
			case "Other Income Expense":
				financialData.OtherIncomeExpense = valueText
			case "Pretax Income":
				financialData.PretaxIncome = valueText
				// ... (add cases for other relevant row titles)
			}
		})
	})

	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/financials", symbol)
	err := c.Visit(url)
	if err != nil {
		log.Printf("Error visiting %s: %v", url, err)
	}

	// Log the final scraped data
	log.Printf("Scraped data for %s: %+v", symbol, financialData)

	return financialData
}

func updateData() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		mu.Lock()
		for symbol := range data {
			data[symbol] = scrapeFinancialData(symbol)
		}
		mu.Unlock()
	}
}
