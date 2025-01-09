package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"urlshortner/middleware"

	_ "github.com/lib/pq"
)

type PageData struct {
	Name string
	URL  string
}

type URLS struct {
	id    int
	url   string
	short string
}

func main() {

	// Connect to the database
	connectStr := "postgresql://postgres:sana12345@localhost:5433/urlshortener?sslmode=disable"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatalln(err)
	}

	// Get all the templates
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	// Create the server
	router := http.NewServeMux()

	// Create the routes for handling request
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", PageData{
			Name: "",
		})
	})
	router.HandleFunc("POST /short", func(w http.ResponseWriter, r *http.Request) { urlForm(w, r, db) })
	router.HandleFunc("GET /{url}", func(w http.ResponseWriter, r *http.Request) { shortURL(w, r, db) })

	// Start the server and log
	server := http.Server{
		Addr:    ":8000",
		Handler: middleware.Logging(router),
	}
	log.Println("Starting server on port :8000")

	err = server.ListenAndServe() // Start the server
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occured:", err)
	}
}

func urlForm(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get all the templates
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	// Parse the form data
	r.ParseForm()

	// Get the data from the submitted form
	urlToShorten := r.FormValue("url")

	// Connect to the database and check if the URL already exists
	// Check if the URL is already in the database
	var res URLS
	var urls []URLS
	var newShortURL string
	var letterCount int
	rows, err := db.Query("SELECT * FROM urls")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res.id, &res.url, &res.short)
		letterCount++
		urls = append(urls, res)

		// Check if the URL is in the database
		if urlToShorten == res.url {
			newShortURL = res.short // Add the already existing short URL
		}
	}

	// Check if the new short url is empty
	if newShortURL == "" {

		// Using letters we will now create the short url
		// Amount of letters the url will have
		_tmpURL := ""
		for j := 0; j <= letterCount/26; j++ {
			newShortURL += _tmpURL

			// Increment the letters from a-z
			for i := 0; i < 26; i++ {
				_tmpURL = string(rune(97 + i))

				// Check if the short URL exists
				good := true
				for _, r := range urls {
					if r.short == newShortURL+_tmpURL || newShortURL+_tmpURL == "" {
						good = false
						break
					}
				}

				// Check if we can use this URL
				if good {
					newShortURL += _tmpURL

					// Add the URL to the database
					_, err := db.Exec("INSERT INTO urls (id, url, short) VALUES ($1, $2, $3)", urls[len(urls)-1].id+1, urlToShorten, newShortURL)
					if err != nil {
						log.Fatalln(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					break
				}
			}
		}
	}

	// Choose the correct schema
	var schema string
	if r.TLS != nil {
		schema = "https://"
	} else {
		schema = "http://"
	}

	// Render the template to show the user their shorten URL
	tmpl.ExecuteTemplate(w, "short.html", PageData{
		URL: schema + r.Host + "/" + newShortURL,
	})
}

func shortURL(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	shortURL := r.PathValue("url")

	// Check if the URL is already in the database
	var res URLS
	var redURL string
	rows, err := db.Query("SELECT * FROM urls")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res.id, &res.url, &res.short)
		if res.short == shortURL {
			redURL = res.url
			break
		}
	}

	http.Redirect(w, r, redURL, http.StatusTemporaryRedirect)
}
