
package main

import (
	"encoding/json"
	"log"
	"net/http"
)






// func respondWithError(w http.ResponseWriter, code int, msg string) {
//     w.WriteHeader(code)
//     json.NewEncoder(w).Encode(errorResponse{Error: msg})
// }

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
//     w.Header().Set("Content-Type", "application/json")
// 	dat, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}
	
	
// 	w.WriteHeader(code)
//     w.Write(dat)
// }
func respondWithError(w http.ResponseWriter, code int, msg string) {
    // w.Header().Set("Content-Type", "application/json")
    // resp := errorResponse{Error: msg}
    // dat, err := json.Marshal(resp)
    // if err != nil {
    //     log.Printf("Error marshalling error JSON: %s", err)
    //     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    //     return
    // }
    // w.WriteHeader(code)
    // w.Write(dat)
    type errorResponse struct {
        Error string `json:"error"`
    }
    
    respondWithJSON(w, code, errorResponse{
        Error: msg,
    })
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error":"Internal Server Error"}`))
        return
        
    }
    w.WriteHeader(code)
    if _, err := w.Write(dat); err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

