package main

import (
	
	"net/http"
	
	"fmt"
)

//middleware method to increment the FileServerHits

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//increment counter
		cfg.fileserverHits.Add(1)

		// continue to next handler
		next.ServeHTTP(w, r)


	})
}


func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html:= fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`,cfg.fileserverHits.Load())
	w.Write([]byte(html))
	
}