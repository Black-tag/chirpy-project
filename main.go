package main 

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

//middleware method to increment the FileServerHits

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//increment counter
		cfg.fileserverHits.Add(1)

		// continue to next handler
		next.ServeHTTP(w, r)


	})
}

// handler to display metrics
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,"Hits: %d", cfg.fileserverHits.Load())
}
//  handler to reset metrics 
func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

// Add the Raediness endpoint or function
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// set content type
	w.Header().Set("Content-Type", "application/json" )

	// Set Status code 
	w.WriteHeader(http.StatusOK)  // return 200

	response := []byte(`{"status": "ready", "message": "server is ready"}`)

	// write response
	_, err := w.Write(response)
	if err != nil {
		fmt.Printf("Error writing response: %v", err)
	}

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader((http.StatusMethodNotAllowed))
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}




func main() {

	cfg := apiConfig{}

	// created a new  mux
	mux := http.NewServeMux()

	//register readiness endpoint
	
	
	mux.HandleFunc("GET /metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /reset", cfg.resetMetrics)
	mux.HandleFunc("/ready", readinessHandler)
	mux.HandleFunc("GET /healthz", healthzHandler)

	// set up fileServer
	mainFs := http.FileServer(http.Dir(".")) 
	assetFs := http.FileServer(http.Dir("./assets"))

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", mainFs)))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", assetFs))

	// creating a http.Server struct 
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	


	fmt.Println("🌐 starting the server on: http://localhost:8080...")
	err:= server.ListenAndServe()
	if err != nil {
		fmt.Printf("🛑 Server Failed: %v\n", err)
	}
	fmt.Println("Server runnig ...")
	



}