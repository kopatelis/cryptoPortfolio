package gecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GeckoClient struct {
	apiKey string
}

func NewClient(apiKey string) *GeckoClient {
	return &GeckoClient{apiKey: apiKey}
}

func (g GeckoClient) Price(symbol string) (Coin, error) {
	client := http.Client{}
	uri := "https://api.coingecko.com/api/v3/simple/price"
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return Coin{}, fmt.Errorf("error creating request: %w", err)
	}

	params := url.Values{}
	params.Add("vs_currencies", "usd")
	params.Add("symbols", symbol)
	params.Add("x_cg_demo_api_key", g.apiKey)

	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return Coin{}, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coin{}, fmt.Errorf("error fail get price: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Coin{}, fmt.Errorf("error reading response body: %w", err)
	}

	return g.ParseCoinPrice(data)
}

func (g GeckoClient) ParseCoinPrice(data []byte) (Coin, error) {
	var PriceResponse map[string]map[string]float64
	err := json.Unmarshal(data, &PriceResponse)
	if err != nil {
		return Coin{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	result := Coin{}

	for coinName, coinPrice := range PriceResponse {
		result.Name = coinName
		result.Price = coinPrice["usd"]
	}

	if len(result.Name) == 0 {
		return Coin{}, fmt.Errorf("no name in response")
	}

	return result, nil
}
