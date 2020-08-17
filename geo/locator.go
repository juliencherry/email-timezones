package geo

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

const (
	utc string = "UTC"
)

type Locator struct {
	IP string

	city string
	timezone string
}

type IPGeolocation struct {
	Timezone	IPGeolocationTimezone `json:"time_zone"`
}

type IPGeolocationTimezone struct {
	Name		string  `json:"name"`
	Offset	int  	  `json:"offset"`
}

func NewLocator(ip string) Locator {
	return Locator{ip, "", ""}
}

func (l Locator) City() string {
	if l.city == "" {
		timezoneParts := strings.SplitAfter(l.Timezone(), "/")
		l.city = strings.Replace(timezoneParts[len(timezoneParts) - 1], "_", " ", -1)
	}

	return l.city

}

func (l Locator) Timezone() string {
	if l.timezone == "" {
		l.timezone = locateTimezone(l.IP)
	}

	return l.timezone;
}

func locateTimezone(ip string) string {

	apiKey := os.Getenv("GEOLOCATION_API_KEY")

	resp, err := http.Get("https://api.ipgeolocation.io/ipgeo?apiKey=" + apiKey + "&ip=" + ip)
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
