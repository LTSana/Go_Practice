package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/tabwriter"

	"golang.org/x/net/html"
)

// Data structor for storing the link data
type Link struct {
	page   string // Full page link
	link   string // Only the link
	status int    // Status code the link returned
}

func scraper(s string) {

	// For goroutine
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Initialize the variables we are going to need
	allLinks := []Link{} // This is to store all the links in a slice
	mainDomain := s      // The main domain to scrape through

	wg.Add(1)
	go diver(&allLinks, mainDomain, "/", &wg, &mu) // Function for crawling the links

	// For adding tabs
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "STATUS\tPAGE\tLINK")

	// Wait for the goroutines to finish
	wg.Wait()

	// Get the concurrencies
	for i := range allLinks {

		// Check if any of the links are dead
		if allLinks[i].status >= 400 {
			fmt.Fprintln(w, "DEAD\t"+allLinks[i].page+"\t"+allLinks[i].link)
		} else {
			fmt.Fprintln(w, "ALIVE\t"+allLinks[i].page+"\t"+allLinks[i].link)
		}
	}

	// Prints out the data to the console in a tab spaced way
	w.Flush()
}

func diver(allLinks *[]Link, mainDomain string, link string, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()

	mu.Lock() // Lock the mutex before accessing the slice
	defer mu.Unlock()

	// Check if the domain is in the link
	link = strings.Replace(link, mainDomain, "", -1)

	// Check if the link is already in the list
	for index := range *allLinks {
		if (*allLinks)[index].link == link || (*allLinks)[index].link == mainDomain+link {
			return nil
		}
	}

	fmt.Print("Checking link: ")
	fmt.Print(mainDomain)
	fmt.Println(link)

	// Check if the link points to a different domain
	if !strings.Contains(link, "http") && !strings.Contains(link, "@") {

		// Make the request to the link
		resp, err := http.Get(mainDomain + link)
		if err != nil {
			fmt.Printf("Error: %v", err)

			// Add the link to the list
			// *allLinks = append(*allLinks, Link{page: mainDomain + link, link: link, status: resp.StatusCode})

			return err
		}
		defer resp.Body.Close()

		// Check if it's a redirect
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			fmt.Println("Redirect...")
			fmt.Println(resp.Header)
		}

		// Add the link to the list
		*allLinks = append(*allLinks, Link{page: mainDomain + link, link: link, status: resp.StatusCode})

		// Tokenize the html string so we can iterage through the html nodes
		z := html.NewTokenizer(resp.Body)
		for {
			tt := z.Next()
			if tt == html.ErrorToken {
				break
			}
			if tt == html.StartTagToken {
				t := z.Token()

				// Check if we found a link tag
				if t.Data == "a" {
					for _, a := range t.Attr {
						if a.Key == "href" {

							// Add another goroutine
							wg.Add(1)

							// Extract the link and crawl it as well
							go diver(allLinks, mainDomain, a.Val, wg, mu)
						}
					}
				}
			}
		}
	} else {
		// Add the link to the list
		*allLinks = append(*allLinks, Link{page: mainDomain + link, link: link, status: 200})
	}

	// Return nil to signal successful
	return nil
}
