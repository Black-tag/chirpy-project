package main

import (
	
	"net/http"
	"fmt"
)



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