package main

import (
	"log"
	"net/http"
)

func main() {
	mainPathParamExample()
}

// Example: 1
func mainPathParamExample() {
	router := http.NewServeMux()
	router.HandleFunc("/item/{id}", pathParamExample)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Println("Server listening on port 8000")
	server.ListenAndServe()
}

func pathParamExample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Item ID: " + id))
}
