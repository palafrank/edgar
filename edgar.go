package edgar

import (
	"io"
	"time"
)

// FilingType is the type definition of various filing types
type FilingType string

// FilingType10Q is a 10-Q quarterly filing of a company with the SEC
const FilingType10Q FilingType = "10-Q"

// FilingType10K is a 10-K annual filing of a company with the SEC
const FilingType10K FilingType = "10-K"

// Filing interface for fetching financial data from a collected filing
type Filing interface {
	Ticker() string
	FiledOn() time.Time
	Type() (FilingType, error)
	ShareCount() (float64, error)
	Revenue() (float64, error)
	CostOfRevenue() (float64, error)
	GrossMargin() (float64, error)
	OperatingIncome() (float64, error)
	OperatingExpense() (float64, error)
	NetIncome() (float64, error)
	TotalEquity() (float64, error)
	ShortTermDebt() (float64, error)
	LongTermDebt() (float64, error)
	CurrentLiabilities() (float64, error)
	CurrentAssets() (float64, error)
	DeferredRevenue() (float64, error)
	RetainedEarnings() (float64, error)
	OperatingCashFlow() (float64, error)
	CapitalExpenditure() (float64, error)
	Dividend() (float64, error)
	WAShares() (float64, error)
	DividendPerShare() (float64, error)
	Interest() (float64, error)
	Cash() (float64, error)
	Securities() (float64, error)
	Goodwill() (float64, error)
	Intangibles() (float64, error)
	CollectedData() []string
}

// CompanyFolder interface used to get filing information about a company
type CompanyFolder interface {

	// Ticker gets the ticker of this company
	Ticker() string

	// CIK gets the CIK assigned to the company
	CIK() string

	// AvailableFilings gets the list of dates of available filings
	AvailableFilings(FilingType) []time.Time

	// Filing gets a filing given a filing type and date of filing.
	Filing(FilingType, time.Time) (Filing, error)

	// Filings gets a list of filings. Parallel fetch.
	Filings(FilingType, ...time.Time) ([]Filing, error)

	// SaveFolder persists the data from the company folder into a writer
	// provided by the user. This stored info can be presented back to
	// the fetcher (using CreateFolder API in fetcher) to recreate the
	// company folder with already parsed data
	SaveFolder(w io.Writer) error

	// String is a dump routine to view the contents of the folder
	String() string

	// Returns a HTML tabulated form of the data
	HTML(FilingType) string
}

// FilingFetcher fetches the filing requested
type FilingFetcher interface {

	// CompanyFolder creates a folder for the company with a list of
	// available filings. No financial data is pulled and the user
	// of the interface can selectively pull financial data into the
	// folder using the CompanyFolder interface
	CompanyFolder(string, ...FilingType) (CompanyFolder, error)

	// CreateFolder creates a company folder using a Reader
	// User can provoder a store of edgar data previous stored
	// by this package (using the Store function of the Company Folder)
	// This function is used to avoid reparsing edgar data and reusing
	// already parsed and stored information.
	CreateFolder(io.Reader, ...FilingType) (CompanyFolder, error)
}
