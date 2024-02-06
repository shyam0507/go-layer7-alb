package main

import (
	"io"
	"net/http"

	"github.com/shyam0507/go-layer7-alb/utility"
	"go.uber.org/zap"
)

type LoadBalancer struct {
	services []string
	client   *http.Client
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	utility.Logger.Info("New Request", zap.String("path", r.URL.Path))
	req, err := http.NewRequest(r.Method, lb.services[0], r.Body)

	if err != nil {
		utility.Logger.Error("Error creating request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err_request := lb.client.Do(req)

	if err_request != nil {
		utility.Logger.Error("Error for request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utility.Logger.Info("Request Resp", zap.Int("status", res.StatusCode))

	defer res.Body.Close()
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func (lb *LoadBalancer) RegisterService(w http.ResponseWriter, r *http.Request) {
}

func (lb *LoadBalancer) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	utility.InitLogger()

	startHTTPServer()
}

func startHTTPServer() {

	loadBalancer := &LoadBalancer{services: []string{"http://localhost:8081/todo"}, client: http.DefaultClient}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", loadBalancer.HandleHealthCheck)
	mux.HandleFunc("/", loadBalancer.HandleRequest)

	http.ListenAndServe(":8080", mux)
}
