package electricity

import "fmt"

func Generate() {
	prices := getPrices()
	fmt.Println("Length", len(prices))
	for _, record := range prices {
		fmt.Printf("%+v", record)
		fmt.Printf(" %+v", (record.FeeDKK + record.SpotPriceDKK))
		fmt.Printf("\n")
	}

	bars := make([]float64, len(prices))
	for i, record := range prices {
		bars[i] = record.SpotPriceDKK + record.FeeDKK
	}
	drawBars(bars, 250, 122)
}
