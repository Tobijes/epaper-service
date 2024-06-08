package electricity

import (
	"image"
	"time"
)

var location, _ = time.LoadLocation("Europe/Copenhagen")

func filterPricesInHours(prices []PriceRecord, hours int) (relevantPrices []PriceRecord) {

	for _, record := range prices {
		now := time.Now().UTC()
		beforeTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		beforeTime = beforeTime.Add(time.Hour * time.Duration(hours))
		if record.StartTimeUTC.Before(beforeTime) {
			relevantPrices = append(relevantPrices, record)
		}
	}
	return
}

func Generate() image.Image {
	prices := getPrices()
	// fmt.Println("Length", len(prices))
	// for _, record := range prices {
	// 	fmt.Printf("%+v", record)
	// 	fmt.Printf(" %+v", (record.FeeDKK + record.SpotPriceDKK))
	// 	fmt.Printf("\n")
	// }

	// Filter to only show next 12 hours
	prices = filterPricesInHours(prices, 12)

	img := drawBars(prices, 250, 122)
	return img
}
