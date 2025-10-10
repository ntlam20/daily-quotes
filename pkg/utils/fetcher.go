package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchQuote() (Quote, error) {
	resp, err := http.Get("https://zenquotes.io/api/random")
	if err != nil {
		return Quote{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data []Quote
	if err := json.Unmarshal(body, &data); err != nil {
		return Quote{}, err
	}

	if len(data) == 0 {
		return Quote{}, fmt.Errorf("no quote returned")
	}

	return data[0], nil
}
