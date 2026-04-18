package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sinde.ru/db"
)

const ipInfoTTL = 24 * time.Hour

type IPInfo struct {
	IP            string  `json:"ip"`
	Success       bool    `json:"success"`
	Type          string  `json:"type"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEU          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`
	Connection    struct {
		ASN    int64  `json:"asn"`
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		ID     string `json:"id"`
		Abbr   string `json:"abbr"`
		IsDST  bool   `json:"is_dst"`
		Offset int64  `json:"offset"`
		UTC    string `json:"utc"`
	} `json:"timezone"`
}

const ipInfoKey = "ipinfo:"

func normalizeIP(ip string) string {
	return ip
}

func ipInfoCacheKey(ip string) string {
	return ipInfoKey + normalizeIP(ip)
}

func getIPInfoFromCache(ctx context.Context, ip string) (*IPInfo, error) {
	s, err := db.RDB.Get(ctx, ipInfoCacheKey(ip)).Result()
	if err != nil || s == "" {
		return nil, err
	}
	var out IPInfo
	if err = json.Unmarshal([]byte(s), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func setIPInfoCache(ctx context.Context, ip string, v *IPInfo) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return db.RDB.Set(ctx, ipInfoCacheKey(ip), b, ipInfoTTL).Err()
}

func fetchIPInfo(ctx context.Context, ip string) (*IPInfo, error) {
	cl := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://ipwho.is/%s", normalizeIP(ip)), nil)
	if err != nil {
		return nil, err
	}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var out IPInfo
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetIPInfo(ctx context.Context, ip string) (*IPInfo, error) {
	v, err := getIPInfoFromCache(ctx, ip)
	if err == nil && v != nil {
		return v, nil
	}
	out, err := fetchIPInfo(ctx, ip)
	if err != nil {
		return nil, err
	}
	if out != nil && out.Success {
		_ = setIPInfoCache(ctx, ip, out)
	}
	return out, nil
}
