package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type GeoIP struct {
	CountryName string `json:"country"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	TimeZone    string `json:"timezone"`
	ISP         string `json:"isp"`
}

func ExtractIPFromMultiAddr(multiAddr string) string {
	parts := strings.Split(strings.Split(multiAddr, "/")[2], "/")
	return parts[0]
}

func GetGeoIP(ip string) *GeoIP {
	geo := &GeoIP{}
	res, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return geo
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return geo
	}

	err = json.Unmarshal(body, &geo)
	if err != nil {
		return geo
	}

	return geo
}
