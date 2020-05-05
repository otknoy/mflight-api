package mflight

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Table struct {
	Id          int64   `xml:"id,attr"`
	Time        string  `xml:"time"`
	Unixtime    int64   `xml:"unixtime"`
	Temperature float32 `xml:"temp"`
	Humidity    float32 `xml:"humi"`
	Illuminance int16   `xml:"illu"`
}

type Response struct {
	Tables []Table `xml:"table"`
}

func getSensorMonitor() (Response, error) {
	baseUrl := "http://192.168.1.16:56002"
	path := "/SensorMonitorV2.xml"

	qs := url.Values{
		"x-KEY_MOBILE_ID":   []string{"xxxx"},
		"x-KEY_UPDATE_DATE": []string{""},
	}

	url := baseUrl + path + "?" + qs.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return Response{}, nil
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, nil
	}

	response := Response{}
	if err := xml.Unmarshal(byteArray, &response); err != nil {
		return Response{}, nil
	}

	return response, nil
}
