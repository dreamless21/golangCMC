package main

import (
	web_app_go "awesomeProject/web-app-go"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const urlAPI = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"

type JsonFile struct {
	Status JsonStatus
	Data   []JsonData
}

type JsonStatus struct {
	Timestamp string `json:"timestamp"`
}

type JsonData struct {
	Id        int                               `json:"id"`
	Name      string                            `json:"name"`
	Symbol    string                            `json:"symbol"`
	Rank      int                               `json:"cmc_rank"`
	DateAdded time.Time                         `json:"date_added"`
	Quote     map[string]map[string]interface{} `json:"quote"`
}

func GetData() string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", urlAPI, nil)
	req.Header = http.Header{
		"Accepts":           {"application/json"},
		"X-CMC_PRO_API_KEY": {web_app_go.APIKey},
	}
	query := req.URL.Query()
	query.Add("start", "1")
	query.Add("limit", "50")
	query.Add("convert", "USD")
	req.URL.RawQuery = query.Encode()
	res, _ := client.Do(req)
	resBody, _ := io.ReadAll(res.Body)
	strToReturn := fmt.Sprintf("%s", resBody)
	return strToReturn
}

func LoadToDB(toJson string) {
	var myJson = JsonFile{}
	db := web_app_go.OpenConnection()
	defer db.Close()
	_ = json.Unmarshal([]byte(toJson), &myJson)
	currentDate, _ := time.Parse(time.RFC3339, myJson.Status.Timestamp)
	currentUTC := currentDate.Add(3 * time.Hour).Format("2006-01-02")
	for _, value := range myJson.Data {
		price := fmt.Sprintf("%.2f", value.Quote["USD"]["price"].(float64))
		percent := fmt.Sprintf("%.3f", value.Quote["USD"]["percent_change_24h"].(float64))
		sqlStmt, _ := db.Prepare("INSERT INTO crypto_day_price(crypto_id, dt, price_usd, percent_change_day, today_rank) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE crypto_id = VALUES(crypto_id)")
		sqlInfo, _ := db.Prepare("INSERT IGNORE INTO crypto_info(crypto_id, crypto_name, crypto_symbol, date_added) VALUES (?, ?, ?, ?)")
		sqlStmt.Exec(value.Id, currentUTC, price, percent, value.Rank)
		sqlInfo.Exec(value.Id, value.Name, value.Symbol, value.DateAdded)
	}
}
