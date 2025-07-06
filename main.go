package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type chirpResponse struct {
    CleanedBody string `json:"cleaned_body"`
}

type chirprequest struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type validResponse struct {
	Valid bool `json:"valid"`
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
//  handler to reset metrics 
func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

//Add the validate_chirp endpoint
func (cfg *apiConfig)validateChirpHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")

	}

	var req chirprequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}

	// validate chirp 

	err = validateChirp(req.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cleanedBody := filterProfanity(req.Body)
	

	// success response
	respondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleanedBody,
	})
	return
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

// helper functions 
func respondWithError(w http.ResponseWriter, code int, msg string) {
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}




func validateChirp(body string) error {
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        return fmt.Errorf("Chirp is too long")
    }
    return nil
}

func filterProfanity(text string) string {
    // Define profane words (case-insensitive)
    profaneWords := map[string]bool{
        "kerfuffle": true,
        "sharbert":  true,
        "fornax":    true,
    }

    // Split text into words while preserving punctuation as separate tokens
    words := strings.Fields(text)
    
    for i, word := range words {
        lowerWord := strings.ToLower(word)
        // Check if word (without punctuation) is profane
        if profaneWords[lowerWord] {
            words[i] = "****"
        }
    }
    
    return strings.Join(words, " ")
}




func main() {

	cfg := apiConfig{}

	// created a new  mux
	mux := http.NewServeMux()

	//register readiness endpoint
	
	mux.HandleFunc("POST /api/validate_chirp",cfg.validateChirpHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)
	mux.HandleFunc("GET /api/ready", readinessHandler)
	mux.HandleFunc("GET /api/healthz", healthzHandler)

	// set up fileServer
	mainFs := http.FileServer(http.Dir(".")) 
	assetFs := http.FileServer(http.Dir("./assets"))

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", mainFs)))
	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/assets", assetFs)))

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