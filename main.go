package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"runtime/debug"
)

var buildTimestamp = ""
var commit = ""

func recordMetrics() {
	fmt.Printf("%[1]s-%[2]s", buildTimestamp, commit)
}

func main() {
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "drone",
			Name:      "exporter_build_info",
			Help:      "Number of blob storage operations waiting to be processed, partitioned by user and type.",
		},
		[]string{
			"version",
		},
	)
	buildInfo.With(prometheus.Labels{"version": fmt.Sprintf("%s-%s", buildTimestamp, commit)}).Set(1)
	prometheus.MustRegister(buildInfo)
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":2112", nil)
	//prometheus.MustRegister(version.NewCollector("app_name"))
	http.Handle("/metrics", promhttp.Handler())

	var sha = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				fmt.Println(setting)
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}
		return ""
	}()

	//buildData, ok := debug.ReadBuildInfo()
	//if !ok {
	//	log.Fatal("err !ok")
	//}

	fmt.Printf("Output: %#v\n", sha)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
