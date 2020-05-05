package mflight

type Temperature float32
type Humidity float32
type Illuminance int16

type Metrics struct {
	Temperature Temperature `xml:"temp"`
	Humidity    Humidity    `xml:"humi"`
	Illuminance Illuminance `xml:"illu"`
}

func GetMetrics() (Metrics, error) {
	res, err := getSensorMonitor()
	if err != nil {
		return Metrics{}, err
	}

	tables := res.Tables
	table := tables[len(tables)-1]

	m := Metrics{
		Temperature(table.Temperature),
		Humidity(table.Humidity),
		Illuminance(table.Illuminance),
	}

	return m, nil
}
