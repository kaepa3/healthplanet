package healthplanet

import (
	"encoding/json"
	"io"
)

type JsonResponce struct {
	BirthDate string `json:"birth_date"`
	Height    string `json:"height"`
	Sex       string `json:"male"`
	Data      Data   `json:"data"`
}

type Data struct {
	Date    string `json:"date"`
	KeyData string `json:"keydata"`
	Model   string `json:"model"`
	Tag     string `json:"tag"`
}

func ConvertToJson(b io.ReadCloser) (*JsonResponce, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	resp := Responce{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
