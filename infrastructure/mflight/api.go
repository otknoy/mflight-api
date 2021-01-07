package mflight

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type table struct {
	ID          int64   `xml:"id,attr"`
	Time        string  `xml:"time"`
	Unixtime    int64   `xml:"unixtime"`
	Temperature float32 `xml:"temp"`
	Humidity    float32 `xml:"humi"`
	Illuminance int16   `xml:"illu"`
}

type response struct {
	Tables []table `xml:"table"`
}

func getSensorMonitor(baseURL, mobileID string) (response, error) {
	url := buildURL(baseURL, mobileID)

	resp, err := http.Get(url)

	if err != nil {
		return response{}, nil
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{}, nil
	}

	res := response{}
	if err := xml.Unmarshal(byteArray, &res); err != nil {
		return response{}, nil
	}

	return res, nil
}

func buildURL(baseURL, mobileID string) string {
	qs := url.Values{
		"x-KEY_MOBILE_ID":   []string{mobileID},
		"x-KEY_UPDATE_DATE": []string{""},
	}

	url, _ := url.Parse(baseURL)
	url.Path = "/SensorMonitorV2.xml"
	url.RawQuery = qs.Encode()

	return url.String()
}
