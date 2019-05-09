package sn

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type SpotterRequestBody struct {
	AppID      string `json:"id"`
	SpotterIds []int  `json:"markers"`
}

type SpotterResponseBody struct {
	Spotters []Spotter `json:"positions"`
}

type Spotter struct {
	LastReport string  `json:"report_at,omitempty"`
	Lat        float64 `json:"lat,string,omitempty"`
	Lon        float64 `json:"lon,string,omitempty"`
	Elev       int     `json:"elev,string,omitempty"`
	Dir        int     `json:"dir,string,omitempty"`
	Gps        bool    `json:"gps,omitempty"`
	Callsign   string  `json:"callsign,omitempty"`
	Email      string  `json:"email,omitempty"`
	Phone      string  `json:"phone,omitempty"`
	Ham        string  `json:"ham,omitempty"`
	HamShow    bool    `json:"ham_show,omitempty"`
	Freq       string  `json:"freq,omitempty"`
	Note       string  `json:"note,omitempty"`
	IM         string  `json:"im,omitempty"`
	Twitter    string  `json:"twitter,omitempty"`
	Web        string  `json:"web,omitempty"`
	Unix       int64   `json:"unix,string,omitempty"`
	First      string  `json:"first,omitempty"`
	Last       string  `json:"last,omitempty"`
	Marker     string  `json:"marker,omitempty"`
}

const snAPIURL = "https://www.spotternetwork.org"

func GetSpotters(appID string, spotterIds []int) ([]Spotter, error) {
	srBody := SpotterRequestBody{
		AppID:      appID,
		SpotterIds: spotterIds,
	}
	data, err := json.Marshal(srBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", snAPIURL+"/positions", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := struct {
		Spotters []Spotter `json:"positions"`
	}{}

	err = json.Unmarshal(body, &tmp)

	return tmp.Spotters, nil
}
