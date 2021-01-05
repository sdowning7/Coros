package main

import (
	"fmt"
	"os"
	"net/http"
	"net/url"
	"golang.org/x/net/html"
)

// an error to represent that we were not able to parse an HTML doc at the given link
type BadLinkError string

func (e BadLinkError) Error() string {
	return fmt.Sprintf("Unable to parse HTML from %v", string(e))
}

type Scraper interface {
	ListURLs() ([]string, error)
	ListExternalURLs() ([]string, error)
}

// a simple web scraper, following the Scraper and Stringer interfaces
type SimpleScraper struct {
	URL string
}

// returns the desired output to be printed according to the Stringer interface
func (s SimpleScraper) String() string {
	urls, err := s.ListExternalURLs()
	if err != nil {
		return fmt.Sprintf("%v %v", s.URL, -1)
	}
	return fmt.Sprintf("%v %v", s.URL, len(urls))
}

//a helper function to tokenize the HTML and return all the links we find
func getAllLinks(resp *http.Response) []string {
	urls := make([]string, 0)

	tokenizer := html.NewTokenizer(resp.Body)

	//This got very ugly very quickly, there might be something better to use
	//iterate through all the tokens produced from the body
	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			//there are no more tokens
			return urls
		} else if tokenType == html.StartTagToken {
			//check to see it it's a link
			token := tokenizer.Token()
			if token.Data == "a" {
				//find the href attribute
				for _, attribute := range token.Attr {
					//append it to our list if it is href
					if attribute.Key == "href" {
						urls = append(urls, attribute.Val)
					}
				}
			}
		}
	}
}

//a helper function to remove all duplicates in a list of strings
func removeDuplicates(strings []string) []string {
	exists := make(map[string]bool)
	noDupes := make([]string, 0)

	for _, str := range strings {
		//if the value does not exist add it to the list and say it exists
		if _, e := exists[str]; !e {
			exists[str] = true
			noDupes = append(noDupes, str)
		}
	}
	return noDupes
}

// return a list of all unique URLs contained in the HTML document pointed to by s.URL
func (s SimpleScraper) ListURLs() ([]string, error) {
	// make a get request for the webpage
	resp, err := http.Get(s.URL)
	if err != nil {
		return make([]string, 0), BadLinkError(s.URL)
	}
	defer resp.Body.Close()

	// extract links from the message content
	urls := getAllLinks(resp)

	//make sure the urls are unique
	urls = removeDuplicates(urls)

	return urls, nil
	
}

// return a list of all unique external URLs contained in the HTML document pointed to by s.URL
func (s SimpleScraper) ListExternalURLs() ([]string, error) {
	urls, err := s.ListURLs()
	externals := make([]string, 0)
	//find all external urls
	for _, v:= range urls {
		link, err := url.Parse(v)
		if link.IsAbs() && err == nil{
			externals = append(externals, v)
		}
	}
	return externals, err
}

// the entry point to our scraper program, parses command line input, then creates and runs the appropriate number of scrapers
func main() {
	
	urls := os.Args[1:]
	for _, v := range urls {
		fmt.Println(SimpleScraper{v})
	}
	
}