package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// https://data.europa.eu/data/datasets/trafikhastigheder?locale=da
// https://wfs-kbhkort.kk.dk/k101/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=k101:trafikhastigheder&outputFormat=json&SRSNAME=EPSG:4326

func main() {
	fmt.Println("Running fetcher")
	client := resty.New()

	params := map[string]string{
		"service":      "WFS",
		"version":      "1.0.0",
		"request":      "GetFeature",
		"typeName":     "k101:trafikhastigheder",
		"outputFormat": "json",
		"SRSNAME":      "EPSG:4326",
	}

	var data Container
	resp, err := client.R().SetQueryParams(params).SetResult(&data).SetHeader("Accept", "application/json").Get("https://wfs-kbhkort.kk.dk/k101/ows")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Fatal(resp.Status())
	}

	fmt.Println(data.TotalFeatures)
}

type Container struct {
	Type          string    `json:"type"`
	TotalFeatures int64     `json:"totalFeatures"`
	Features      []Feature `json:"features"`
}

type Feature struct {
	Type         string     `json:"type"`
	ID           string     `json:"id"`
	Geometry     Geometry   `json:"geometry"`
	GeometryName string     `json:"geometry_name"`
	Properties   Properties `json:"properties"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type Properties struct {
	ID                    int     `json:"id"`
	City                  string  `json:"bydel"`
	RoadID                string  `json:"vejid"`
	Roadname              string  `json:"vejnavn"`
	FromStation           float64 `json:"frastation"`
	ToStation             float64 `json:"tilstation"`
	Speedlimit            int     `json:"hastighedsgraense"`
	RecommendedSpeedlimit *int    `json:"anb_hastighedsgraense"`
}
