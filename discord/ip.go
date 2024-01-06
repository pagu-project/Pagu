package discord

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type GeoIP struct {
	CountryName string `json:"country"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	TimeZone    string `json:"timezone"`
	ISP         string `json:"isp"`
}

func getGeoIP(ip string) *GeoIP {
	geo := &GeoIP{}
	req, err := http.NewRequest("GET", "http://ip-api.com/json/"+ip, nil)
	if err != nil {
		return geo
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Millisecond)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return geo
	}
	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return geo
	}

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &geo)
	if err != nil {
		return geo
	}

	return geo
}
