package cmc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/synternet/price-publisher/internal/config"
)

type CmcQuotesLatest struct {
	Status Status              `json:"status"`
	Data   map[string]DataItem `json:"data"`
}

type Status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage *string   `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
	Notice       *string   `json:"notice"`
}

type DataItem struct {
	ID                            int              `json:"id"`
	Name                          string           `json:"name"`
	Symbol                        string           `json:"symbol"`
	Slug                          string           `json:"slug"`
	NumMarketPairs                int              `json:"num_market_pairs"`
	DateAdded                     time.Time        `json:"date_added"`
	Tags                          []Tag            `json:"tags"`
	MaxSupply                     *float64         `json:"max_supply"`
	CirculatingSupply             float64          `json:"circulating_supply"`
	TotalSupply                   float64          `json:"total_supply"`
	Platform                      *Platform        `json:"platform"`
	IsActive                      int              `json:"is_active"`
	InfiniteSupply                bool             `json:"infinite_supply"`
	CmcRank                       int              `json:"cmc_rank"`
	IsFiat                        int              `json:"is_fiat"`
	SelfReportedCirculatingSupply *float64         `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         *float64         `json:"self_reported_market_cap"`
	TvlRatio                      *float64         `json:"tvl_ratio"`
	LastUpdated                   time.Time        `json:"last_updated"`
	Quote                         map[string]Quote `json:"quote"`
}

type Tag struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

type Quote struct {
	Price                 float64   `json:"price"`
	Volume24h             float64   `json:"volume_24h"`
	VolumeChange24h       float64   `json:"volume_change_24h"`
	PercentChange1h       float64   `json:"percent_change_1h"`
	PercentChange24h      float64   `json:"percent_change_24h"`
	PercentChange7d       float64   `json:"percent_change_7d"`
	PercentChange30d      float64   `json:"percent_change_30d"`
	PercentChange60d      float64   `json:"percent_change_60d"`
	PercentChange90d      float64   `json:"percent_change_90d"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapDominance    float64   `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap"`
	Tvl                   *float64  `json:"tvl"`
	LastUpdated           time.Time `json:"last_updated"`
}

type QuoteInfo struct {
	Price           float64 `json:"price"`
	PercentChange24 float64 `json:"price_percent_change_24h"`
	LastUpdated     int64   `json:"last_updated"`
}

func RetrievePrices(cfg config.CmcConfig) (*CmcQuotesLatest, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?id=%s", cfg.Ids), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CMC_PRO_API_KEY", cfg.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-2xx status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cmcQuotesLatest CmcQuotesLatest
	err = json.Unmarshal(body, &cmcQuotesLatest)
	if err != nil {
		return nil, err
	}

	return &cmcQuotesLatest, nil
}
