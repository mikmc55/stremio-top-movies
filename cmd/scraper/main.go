package main

import (
	"flag"
	"strings"
)

var (
	dataDir = flag.String("dataDir", ".", "Location of the data directory. This is where the CSV files will be written.")
)

func main() {
	flag.Parse()

	// Clean input
	if strings.HasSuffix(*dataDir, "/") {
		*dataDir = strings.TrimRight(*dataDir, "/")
	}

	imdbClient := newIMDbClient()
	rtClient := newRTClient(imdbClient) // Changed function name to camelCase

	imdbClient.scrapeTop250(getFilePath("imdb-top-250.csv"))
	imdbClient.scrapeMostPopular(getFilePath("imdb-most-popular.csv"))
	imdbClient.scrapeBoxOfficeUSWeekend(getFilePath("top-box-office-us.csv"))
	rtClient.scrapeCertifiedFreshDVDstreaming(getFilePath("rt-certified-fresh.csv"))
}

// getFilePath constructs the full file path for a given filename
func getFilePath(filename string) string {
	return *dataDir + "/" + filename
}
