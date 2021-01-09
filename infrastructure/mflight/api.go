package mflight

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Table struct {
	ID          int64   `xml:"id,attr"`
	Time        string  `xml:"time"`
	Unixtime    int64   `xml:"unixtime"`
	Temperature float32 `xml:"temp"`
	Humidity    float32 `xml:"humi"`
	Illuminance int16   `xml:"illu"`
}

type Response struct {
	Tables []Table `xml:"table"`
}

type MfLightClient interface {
	GetSensorMonitor() (*Response, error)
}

func NewMfLightClient(baseURL, mobileID string) MfLightClient {
	return &mfLightClient{baseURL, mobileID}
}

type mfLightClient struct {
	baseURL  string
	mobileID string
}

func (c *mfLightClient) GetSensorMonitor() (*Response, error) {
	url := buildURL(c.baseURL, c.mobileID)

	resp, err := http.Get(url)

	if err != nil {
		return &Response{}, nil
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{}, nil
	}

	res := &Response{}
	if err := xml.Unmarshal(byteArray, &res); err != nil {
		return &Response{}, nil
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
