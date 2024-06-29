package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type IMDbClient struct {
	httpClient *http.Client
}

// Creates a new IMDbClient with a timeout set for HTTP requests
func newIMDbClient() *IMDbClient {
	return &IMDbClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Scrapes the IMDb Top 250 movies and saves to a CSV file
func (c *IMDbClient) scrapeTop250(filePath string) {
	req, err := http.NewRequest("GET", "https://www.imdb.com/chart/top/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept-language", "en-US")

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	// Writing header
	header := []string{"rank", "title", "year", "IMDb ID"}
	if err := csvWriter.Write(header); err != nil {
		log.Fatal(err)
	}

	doc.Find(".lister-list tr").Each(func(i int, s *goquery.Selection) {
		rank := i + 1
		title := s.Find(".titleColumn a").Text()
		href, _ := s.Find(".titleColumn a").Attr("href")
		year := s.Find(".titleColumn span").Text()
		year = strings.Trim(year, "()")
		id := strings.Split(href, "/")[2]

		record := []string{strconv.Itoa(rank), title, year, id}
		if err := csvWriter.Write(record); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v. %v (%v); ID: %v\n", rank, title, year, id)
	})
}

// Scrapes the IMDb Most Popular movies and saves to a CSV file
func (c *IMDbClient) scrapeMostPopular(filePath string) {
	req, err := http.NewRequest("GET", "https://www.imdb.com/chart/moviemeter", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept-language", "en-US")

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	// Writing header
	header := []string{"rank", "title", "year", "IMDb ID"}
	if err := csvWriter.Write(header); err != nil {
		log.Fatal(err)
	}

	doc.Find(".lister-list tr").Each(func(i int, s *goquery.Selection) {
		rank := i + 1
		title := s.Find(".titleColumn a").Text()
		href, _ := s.Find(".titleColumn a").Attr("href")
		year := s.Find(".titleColumn .secondaryInfo").Text()
		year = strings.Trim(year, "()") // Adjust this line to properly handle the year extraction
		id := strings.Split(href, "/")[2]

		record := []string{strconv.Itoa(rank), title, year, id}
		if err := csvWriter.Write(record); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v. %v (%v); ID: %v\n", rank, title, year, id)
	})
}

// Scrapes the IMDb Box Office US Weekend chart and saves to a CSV file
func (c *IMDbClient) scrapeBoxOfficeUSWeekend(filePath string) {
	req, err := http.NewRequest("GET", "https://www.imdb.com/chart/boxoffice", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept-language", "en-US")

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	// Writing header
	header := []string{"rank", "title", "IMDb ID"}
	if err := csvWriter.Write(header); err != nil {
		log.Fatal(err)
	}

	doc.Find(".chart tbody tr").Each(func(i int, s *goquery.Selection) {
		rank := i + 1
		title := s.Find(".titleColumn a").Text()
		href, _ := s.Find(".titleColumn a").Attr("href")
		id := strings.Split(href, "/")[2]
		id = strings.Split(id, "?")[0] // Ensure ID is clean and not polluted by query parameters

		record := []string{strconv.Itoa(rank), title, id}
		if err := csvWriter.Write(record); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v. %v; ID: %v\n", rank, title, id)
	})
}

// Retrieves the IMDb ID for a given movie title
func (c *IMDbClient) getID(title string) string {
	title = url.QueryEscape(title)
	req, err := http.NewRequest("GET", "https://www.imdb.com/find?q="+title, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept-language", "en-US")

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var id string
	doc.Find(".result_text").Each(func(i int, s *goquery.Selection) {
		if i > 0 { // Only consider the first result
			return
		}
		href, _ := s.Find("a").Attr("href")
		id = strings.Split(href, "/")[2]
	})
	return id
}
