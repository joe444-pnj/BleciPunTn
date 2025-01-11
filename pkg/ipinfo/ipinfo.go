package ipinfo

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Location string `json:"loc"`
	Org      string `json:"org"`
	Timezone string `json:"timezone"`
	Postal   string `json:"postal"`
}

func FetchIPInfo(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	var ipInfo IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return nil, err
	}

	if ipInfo.Hostname == "" {
		hostnames, _ := net.LookupAddr(ip)
		if len(hostnames) > 0 {
			ipInfo.Hostname = strings.Join(hostnames, ", ")
		} else {
			ipInfo.Hostname = "Unknown"
		}
	}

	return &ipInfo, nil
}
