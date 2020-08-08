package geo

import (
	"encoding/json"
	"net/http"
	"os"
)

type Locator struct {
	IP string
}

type IPGeolocation struct {
	Timezone	IPGeolocationTimezone `json:"time_zone"`
}

type IPGeolocationTimezone struct {
	Name		string  `json:"name"`
	Offset	int  	  `json:"offset"`
}

var (
	utc string = "UTC"
)

func (l Locator) Timezone() string {

	apiKey := os.Getenv("GEOLOCATION_API_KEY")

	resp, err := http.Get("https://api.ipgeolocation.io/ipgeo?apiKey=" + apiKey + "&ip=" + l.IP)
	if err != nil {
		return utc
	}
	defer resp.Body.Close()

	var location IPGeolocation
	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		return utc
	}

	if location.Timezone.Name == "" {
		return utc
	}

	return location.Timezone.Name
}
