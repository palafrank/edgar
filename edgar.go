package edgar

import (
	"errors"
	"log"
)

type FilingType string

var (
	//Filing types
	FilingType10Q FilingType = "10-Q"
	FilingType10K FilingType = "10-K"
)

// Date defines an interface for filing date
type Date interface {
	Day() int
	Month() int
	Year() int
	String() string
}

// Filing interface for fetching financial data
type Filing interface {
	Ticker() string
	Year() int
	Month() int
	Type() (FilingType, error)
	ShareCount() (int64, error)
	Revenue() (int64, error)
	CostOfRevenue() (int64, error)
	GrossMargin() (int64, error)
	OperatingIncome() (int64, error)
	OperatingExpense() (int64, error)
	NetIncome() (int64, error)
	TotalEquity() (int64, error)
	ShortTermDebt() (int64, error)
	LongTermDebt() (int64, error)
	CurrentLiabilities() (int64, error)
	DeferredRevenue() (int64, error)
	RetainedEarnings() (int64, error)
	OperatingCashFlow() (int64, error)
	CapitalExpenditure() (int64, error)
}

// Company interface used to get information and filing about a company
type CompanyFolder interface {

	// GetTicker gets the ticker of this company
	Ticker() string

	//AvailableFilings gets the list of keys to the filing available
	AvailableFilings(FilingType) []Date

	// GetFiling gets a filing given a Company.
	Filing(FilingType, Date) (Filing, error)
}

// FilingFetcher fetches the filing requested
type FilingFetcher interface {

	// GetFilings creates a folder for the company with a list of
	// available filings. No financial data is pulled and the user
	// of the interface can selectively pull financial data into the
	// folder using the CompanyFolder interface
	CompanyFolder(string, ...FilingType) (CompanyFolder, error)
}

type fetcher struct {
}

/*
	Sequence of extracting financial data:
	    - Input: Ticker symbol and type of filing
			- Step 1: Using input get the links available for the query
			    - The returned map is indexed on date and contains links to the filing
			- Step 2: For each link
			    - Get the documents related to that filing ex. Entity, Balance Sheet
					- For each document get the relevant information and return the data
					- Collect the data into a report
					- Add the report under the TICKER and the date in that order

*/
func (f *fetcher) CompanyFolder(
	ticker string,
	fileTypes ...FilingType) (CompanyFolder, error) {

	comp := newCompany(ticker)

	for _, t := range fileTypes {
		comp.FilingLinks[t] = getFilingLinks(ticker, t)
	}
	return comp, nil
}

func (f *fetcher) CompanyFiling(
	ticker string,
	fileType FilingType,
	d Date) (Filing, error) {

	filingLinks := getFilingLinks(ticker, fileType)
	for key, val := range filingLinks {
		if key != d.String() {
			continue
		}
		log.Println("Geting filing for", ticker, "filed on", key)
		file := new(filing)
		file.FinData = getFinancialData(val, fileType)
		file.Date = key
		file.Company = ticker
		return file, nil
	}
	return nil, errors.New("Could not find the requested filing")
}

func NewFilingFetcher() FilingFetcher {
	return &fetcher{}
}
