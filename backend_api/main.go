package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func main() {
	backendApi()
}

func example1() {
	// Home page basically
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! NIGGA")
	})

	if err := http.ListenAndServe(":8086", nil); err != nil {
		fmt.Println("Server error: ", err)
	}
}

func example2() {
	// Home page basically
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! NIGGA")
	})

	// Second route is users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the users page!")
	})

	if err := http.ListenAndServe(":8086", nil); err != nil {
		fmt.Println("Server error: ", err)
	}
}

func example3() {

	// Create client
	client := &http.Client{}

	// Make a request to the given URL
	resp, err := client.Get("https://api.npoint.io/472dd0f13f829e1f41dc")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read all the content of the body into a buffer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error: ", err)
		return
	}

	// Print out the response
	fmt.Println("Response: ", string(body))
}

func example4() {

	// Set up a router and check its request method
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			fmt.Fprintf(w, "You made a GET request!")
		case "POST":

			// Get the payload from the request
			r.ParseMultipartForm(1024)
			for key, val := range r.Form {
				fmt.Printf("%s = %s\n", key, val)
			}

			fmt.Printf("username: %s\n", r.FormValue("username"))

			fmt.Fprintf(w, "You made a POST request! EX4")
		default:
			http.Error(w, "Invalid request method! Only GET or POST allowed.", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	if err := http.ListenAndServe(":8086", nil); err != nil {
		fmt.Println("Server error: ", err)
	}
}

func backendApi() {

	// For logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	router := http.NewServeMux()

	// Handle addition request
	router.HandleFunc("POST /add", handleAdd)
	router.HandleFunc("POST /multi", handleMulti)
	router.HandleFunc("POST /div", handleDiv)

	// Start the server
	srv := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	logger.Info("Server listening on port 8000")
	srv.ListenAndServe()
}

func handleAdd(w http.ResponseWriter, r *http.Request) {

	// Prepare the form for reading
	r.ParseMultipartForm(1)

	// Check if the user provided 'a' & 'b' field
	if r.FormValue("a") == "" || r.FormValue("b") == "" {
		http.Error(w, "Please provide 'a' and 'b' fields", http.StatusBadRequest)
		return
	}

	// Get the values and convert them to float 64 with bitsize 64
	_a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
	_b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
	total := fmt.Sprintf("Total: %f", (_a + _b))

	// Return the response
	fmt.Fprintf(w, "%s", total)
}

func handleMulti(w http.ResponseWriter, r *http.Request) {

	// Prepare the form for reading
	r.ParseMultipartForm(1)

	// Check if the user provided 'a' & 'b' field
	if r.FormValue("a") == "" || r.FormValue("b") == "" {
		http.Error(w, "Please provide 'a' and 'b' fields", http.StatusBadRequest)
		return
	}

	// Get the values and convert them to float 64 with bitsize 64
	_a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
	_b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
	total := fmt.Sprintf("Total: %f", (_a * _b))

	// Return the response
	fmt.Fprintf(w, "%s", total)
}

func handleDiv(w http.ResponseWriter, r *http.Request) {

	// Prepare the form for reading
	r.ParseMultipartForm(1)

	// Check if the user provided 'a' & 'b' field
	if r.FormValue("a") == "" || r.FormValue("b") == "" {
		http.Error(w, "Please provide 'a' and 'b' fields", http.StatusBadRequest)
		return
	}

	// Can't divide by zero
	if r.FormValue("b") == "0" {
		http.Error(w, "Can't divide by zero", http.StatusBadRequest)
		return
	}

	// Get the values and convert them to float 64 with bitsize 64
	_a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
	_b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
	total := fmt.Sprintf("Total: %f", (_a / _b))

	// Return the response
	fmt.Fprintf(w, "%s", total)
}
