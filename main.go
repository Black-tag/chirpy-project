package main 

import (
	"fmt"
	"net/http"
)






func main() {

	// created a new  mux
	mux := http.NewServeMux()



	// creating a http.Server struct 
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	fmt.Println("starting the server on: 8080...")
	err:= server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server Failed: %v\n", err)
	}



}