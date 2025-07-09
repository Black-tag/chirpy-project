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


	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
	}
	

	
	// defer dbConn.Close()
	// dbQueries := database.New(db)
	// Get platform (default to "dev" if not set)
    
	// created a new  mux
	mux := http.NewServeMux()

	//register readiness endpoint
	
	mux.HandleFunc("POST /api/validate_chirp",cfg.validateChirpHandler)
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