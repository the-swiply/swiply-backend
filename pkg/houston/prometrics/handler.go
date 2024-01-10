package prometrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"net/http"
)

func RunMetricsHandler(path, addr string) {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())

	go func() {
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			loggy.Fatal("can't run prometheus metrics handler: ", err)
		}
	}()

	loggy.Infoln("metrics handler started")
}
