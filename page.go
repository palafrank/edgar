package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	baseURL   string = "https://www.sec.gov/"
	queryURL  string = "cgi-bin/browse-edgar?action=getcompany&CIK=%s&type=%s&dateb=&owner=exclude&count=10"
	searchURL string = baseURL + queryURL
)

func createQueryURL(symbol string, docType filingType) string {
	return fmt.Sprintf(searchURL, symbol, docType)
}

func getPage(url string) io.Reader {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Query to SEC page failed: ", err)
		return nil
	}
	resp.Body.Close()
	return resp.Body
}

func getFilingLinks(ticker string, fileType filingType) []string {
	url := createQueryURL(ticker, fileType)
	resp := getPage(url)
	if resp == nil {
		return nil
	}
	return queryPageParser(resp)

}

func getFilingPage(url string, fileType filingType) map[filingDocType]string {
	url = baseURL + url
	resp := getPage(url)
	if resp == nil {
		return nil
	}
	return filingPageParser(resp, fileType)
}

func getEnityPage(url string) (*EntityData, error) {
	url = baseURL + url
	resp := getPage(url)
	if resp == nil {
		return nil, errors.New("Could not find the Entity page " + url)
	}
	return getEntityData(resp)
}

func getOpsPage(url string) (*OpsData, error) {
	url = baseURL + url
	resp := getPage(url)
	if resp == nil {
		return nil, errors.New("Could not find the Operations page " + url)
	}
	return getOpsData(resp)
}
