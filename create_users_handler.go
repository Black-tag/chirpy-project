
package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"
	"github.com/google/uuid"
	// "context"
)
type createUserRequest struct {
	Email string `json:"email"`
	// Password string `json:"password"`
}

type createUserResponse struct {
        ID        uuid.UUID `json:"id"`
        Email     string    `json:"email"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }


// handler to create users
func (cfg *apiConfig)createUserHandler(w http.ResponseWriter, r *http.Request){
	// method check
	if r.Method != "POST"{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//  if cfg == nil {
    //     http.Error(w, "Server configuration error", http.StatusInternalServerError)
    //     return
    // }

    // // 2. Check if DB is initialized
    // if cfg.DB == nil {
    //     http.Error(w, "Database not initialized", http.StatusInternalServerError)
    //     return
    // }

	// log.Println("incoming requests to api/users")
	
	w.Header().Set("Content-Type", "application/json")
	
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err!=nil {
		respondWithError(w, http.StatusBadRequest,"Invalid Request")
	}
	// validate the input 
	if req.Email == ""  {
		respondWithError(w, http.StatusBadRequest, "Inavalid request")
		return
	}

	// ctx := context.Background()
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("❌ DB error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, createUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}