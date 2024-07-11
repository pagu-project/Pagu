package utils

import (
	"context"
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

func GetGeoIP(ctx context.Context, ip string) *GeoIP {
	geo := &GeoIP{}
	endpoint := "http://ip-api.com/json/" + ip
	cli := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return geo
	}

	resp, err := cli.Do(req)
	if err != nil {
		return geo
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return geo
	}

	err = json.Unmarshal(body, &geo)
	if err != nil {
		return geo
	}

	return geo
}
