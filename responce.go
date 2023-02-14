package healthplanet

import (
	"encoding/json"
	"io"
)

type JsonResponce struct {
	BirthDate string `json:"birth_date"`
	Height    string `json:"height"`
	Sex       string `json:"male"`
	Data      []Data `json:"data"`
}

type Data struct {
	// 測定日
	Date string `json:"date"`
	// 測定データ
	KeyData string `json:"keydata"`
	// 測定機器名
	Model string `json:"model"`
	// 測定部位
	Tag string `json:"tag"`
}

func ConvertToJson(b io.ReadCloser) (*JsonResponce, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	resp := JsonResponce{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
