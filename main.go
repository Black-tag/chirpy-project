package main 

import (
	"fmt"
	"net/http"
)






func main() {

	// created a new  mux
	mux := http.NewServeMux()

	// set up fileServer
	filePath := http.Dir(".")
	handler := http.FileServer(filePath)
	mux.Handle("/", handler)

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