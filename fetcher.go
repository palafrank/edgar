package edgar

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type fetcher struct {
	folders map[string]*company
}

// CompanyFolder creates a new folder and populates it with the filing filing
// links available for the list of filing types
func (f *fetcher) CompanyFolder(
	ticker string,
	fileTypes ...FilingType) (CompanyFolder, error) {

	comp, ok := f.folders[ticker]
	if !ok {
		comp = newCompany(ticker)
		if comp.cik == "" {
			return nil, errors.New("Could not find the CIK for the given ticker")
		}
		f.folders[ticker] = comp
		for _, t := range fileTypes {
			comp.addFilingLinks(t, getFilingLinks(ticker, t))
		}
	}
	return comp, nil
}

// CreateFolder Reads from the reader into a new company folder
func (f *fetcher) CreateFolder(
	r io.Reader,
	fileTypes ...FilingType) (CompanyFolder, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := newCompany("")

	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	// Populate the CIK
	c.cik = getCompanyCIK(c.Ticker())
	if c.cik == "" {
		return nil, errors.New("Could not find the CIK for the given ticker")
	}

	f.folders[c.Ticker()] = c
	// Get all the latest links for all the filing types
	for _, key := range fileTypes {
		c.addFilingLinks(key, getFilingLinks(c.Ticker(), key))
	}
	return c, nil
}

// NewFilingFetcher creates a new empty filing fetcher
func NewFilingFetcher() FilingFetcher {
	return &fetcher{folders: make(map[string]*company)}
}
