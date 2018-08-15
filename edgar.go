package main

type filingType string
type filingDocType string

var (
	filingType10Q filingType = "10-Q"
	filingType10K filingType = "10-K"

	filingDocOps filingDocType = "Operations"
	filingDocInc filingDocType = "Income"
	filingDocBS  filingDocType = "Assets"
	filingDocCF  filingDocType = "Cash Flow"
	filingDocEN  filingDocType = "Entity Info"
	filingDocIg  filingDocType = "Ignore"
)

func main() {

	/*
	   First run the query and get all the links for the filings
	   This will give an array of links to each of the filings.
	   The number of links depends on the query
	*/
	filingLinks := getFilingLinks("AAPL", filingType10Q)

	/*
	   Go through each filing
	   - Get the filing page and find all the links to the associated reports
	         - This will give a map of document types mapped to the corresponding links
	   - Go through each of the entries in the map and for each document type get the data needed
	*/

	for _, link := range filingLinks {
		getFilingPage(link, filingType10Q)
		break
	}

}
