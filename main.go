package main

import (
	
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	

	"github.com/blacktag/chirpy-project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
	
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

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	
	platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatal("Platform must be set")
    }
	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(dbConn)
	secret := os.Getenv("SECRET")
	

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
		secret: secret,
	}
	


    
	// created a new  mux
	mux := http.NewServeMux()

	//register readiness endpoint
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirpByIDHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateCredentialsHandler)
	mux.HandleFunc("/api/revoke", cfg.revokeHandler)
	mux.HandleFunc("/api/refresh", cfg.refreshTokenHandler)
	mux.HandleFunc("POST /api/login", cfg.loginUserHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpByIDHandler)
	mux.HandleFunc("GET /api/chirps",cfg.getChirpHandler)
	mux.HandleFunc("POST /api/chirps",cfg.createChirpHandler)
	mux.HandleFunc("POST /api/users",cfg.createUserHandler)
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
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("🛑 Server Failed: %v\n", err)
	}
	fmt.Println("Server runnig ...")
	



}