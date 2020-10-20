package main

import (
	metering "hypercloud-api-server/metering"
	"hypercloud-api-server/namespace"
	user "hypercloud-api-server/user"
	"net/http"

	"github.com/robfig/cron"
	"k8s.io/klog"
)

func main() {
	// Metering Cron Job
	cronJob := cron.New()
	cronJob.AddFunc("0 */5 * ? * *", metering.MeteringJob)
	cronJob.Start()

	// Req multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/user", serveUser)
	mux.HandleFunc("/metering", serveMetering)
	mux.HandleFunc("/namespace", serveNamespace)
	//mux.HandleFunc("/namespaceClaim", serveNamespaceClaim)

	// HTTP Server Start
	klog.Info("Starting Hypercloud-Operator-API server...")
	if err := http.ListenAndServe(":80", mux); err != nil {
		klog.Errorf("Failed to listen and serve Hypercloud-Operator-API server: %s", err)
	}
	klog.Info("Started Hypercloud-Operator-API server")

}

func serveNamespace(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		namespace.Get(res, req)
	case http.MethodPut:
		namespace.Put(res, req)
	case http.MethodOptions:
		namespace.Options(res, req)
	default:
		//error
	}
}

func serveUser(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		user.Post(res, req)
	case http.MethodDelete:
		user.Delete(res, req)
	case http.MethodOptions:
		user.Options(res, req)
	default:
		//error
	}
}

func serveMetering(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		metering.Get(res, req)
	case http.MethodOptions:
		metering.Options(res, req)
	default:
		//error
	}
}
