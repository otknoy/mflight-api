package main

import (
	"fmt"
	"log"
	"mflight-exporter/config"
	"mflight-exporter/mflight"
	"net/http"
)

const (
	namespace = "MultiFunctionLight"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sensor := mflight.NewMfLight(c.MfLight.URL, c.MfLight.MobileID)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics, _ := sensor.GetMetrics()

		fmt.Fprintf(w, "# HELP %s_temperature multifunction light temperature\n", namespace)
		fmt.Fprintf(w, "# TYPE %s_temperature gauge\n", namespace)
		fmt.Fprintf(w, "%s_temperature %.1f\n", namespace, metrics.Temperature)

		fmt.Fprintf(w, "# HELP %s_humidity multifunction light humidity\n", namespace)
		fmt.Fprintf(w, "# TYPE %s_humidity gauge\n", namespace)
		fmt.Fprintf(w, "%s_humidity %.1f\n", namespace, metrics.Humidity)

		fmt.Fprintf(w, "# HELP %s_illuminance multifunction light illuminance\n", namespace)
		fmt.Fprintf(w, "# TYPE %s_illuminance gauge\n", namespace)
		fmt.Fprintf(w, "%s_illuminance %d\n", namespace, metrics.Illuminance)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
