package httpclient

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is http client interface to get metrics.
type Client interface {
	GetSensorMonitor(ctx context.Context) (*Response, error)
}

// Response is struct to represent root XML element in mflight response.
type Response struct {
	Tables []Table `xml:"table"`
}

// Table is struct to represent <table /> XML element.
type Table struct {
	ID          int64   `xml:"id,attr"`
	Time        string  `xml:"time"`
	Unixtime    int64   `xml:"unixtime"`
	Temperature float32 `xml:"temp"`
	Humidity    float32 `xml:"humi"`
	Illuminance int16   `xml:"illu"`
}

// NewClient creates a new Client.
func NewClient(c *http.Client, baseURL, mobileID string) Client {
	return &client{c, baseURL, mobileID}
}

type client struct {
	client *http.Client

	baseURL  string
	mobileID string
}

// GetSensorMonitor returns mflight response.
func (c *client) GetSensorMonitor(ctx context.Context) (*Response, error) {
	r := buildRequestWithContext(ctx, c.baseURL, c.mobileID)

	resp, err := c.client.Do(r)
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

func buildRequestWithContext(ctx context.Context, baseURL, mobileID string) *http.Request {
	qs := url.Values{
		"x-KEY_MOBILE_ID":   []string{mobileID},
		"x-KEY_UPDATE_DATE": []string{""},
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	req.URL.Path = "/SensorMonitorV2.xml"
	req.URL.RawQuery = qs.Encode()

	return req
}
