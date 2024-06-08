package electricity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Total   int    `json:"total"`
	Filters string `json:"filters"`
	Sort    string `json:"sort"`
	Dataset string `json:"dataset"`
	Records []struct {
		HourUTC      string  `json:"HourUTC"`
		HourDK       string  `json:"HourDK"`
		PriceArea    string  `json:"PriceArea"`
		SpotPriceDKK float64 `json:"SpotPriceDKK"`
		SpotPriceEUR float64 `json:"SpotPriceEUR"`
	} `json:"records"`
}

type PriceRecord struct {
	StartTimeUTC time.Time
	SpotPriceDKK float64
	FeeDKK       float64
}

// https://www.energidataservice.dk/tso-electricity/Elspotprices
func getPrices() (prices []PriceRecord) {
	baseUrl := "https://api.energidataservice.dk/dataset/Elspotprices"
	url, err := url.Parse(baseUrl)
	if err != nil {
		return
	}

	params := url.Query()

	now := time.Now().UTC()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	startTimeFormatted := startTime.Format("2006-01-02T15:04")
	fmt.Println(startTimeFormatted)
	params.Add("start", startTimeFormatted)
	params.Add("filter", "{\"PriceArea\":[\"DK2\"]}")
	params.Add("sort", "HourUTC ASC")

	url.RawQuery = params.Encode()
	fmt.Println(url)
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(responseData))

	var responseObject Response
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%+v\n", responseObject)

	prices = make([]PriceRecord, len(responseObject.Records))
	for i, record := range responseObject.Records {
		t, err := time.Parse("2006-01-02T15:04:05", record.HourUTC)
		if err != nil {
			fmt.Println(err)
			return
		}
		price := record.SpotPriceDKK / 1000.0 // SpotPrice is DKK/MWh
		price = price * 1.25                  // VAT
		fee := computeFee(t)
		prices[i] = PriceRecord{StartTimeUTC: t, SpotPriceDKK: price, FeeDKK: fee}
	}

	return prices
}

// Hour (0-index) to tarif
var localTariffs = map[int]float64{
	0:  0.15187,
	1:  0.15187,
	2:  0.15187,
	3:  0.15187,
	4:  0.15187,
	5:  0.15187,
	6:  0.22775,
	7:  0.22775,
	8:  0.22775,
	9:  0.22775,
	10: 0.22775,
	11: 0.22775,
	12: 0.22775,
	13: 0.22775,
	14: 0.22775,
	15: 0.22775,
	16: 0.22775,
	17: 0.59225,
	18: 0.59225,
	19: 0.59225,
	20: 0.59225,
	21: 0.22775,
	22: 0.22775,
	23: 0.22775,
}

// https://stromligning.dk/live
func computeFee(startTime time.Time) float64 {
	loc, _ := time.LoadLocation("Europe/Copenhagen")
	feeState := 0.95125
	feeSystemTariff := 0.06375
	feeNetTariff := 0.09250
	feeLocalTariff := localTariffs[startTime.In(loc).Hour()]
	return feeState + feeSystemTariff + feeNetTariff + feeLocalTariff
}
