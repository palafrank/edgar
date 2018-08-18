package main

import (
	"log"
)

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

func main() {
	ticker := "AAPL"
	fileType := filingType10K

	var company Company
	/*
		   First run the query and get all the links for the filings of a certain type
			 Return:
			   Map of filing links indexed by date of filing
	*/
	filingLinks := getFilingLinks(ticker, fileType)

	company.Ticker = ticker

	for key, val := range filingLinks {
		filing := new(Filing)
		log.Println("Geting filing for", ticker, "filed on", key)
		filing.FinData = getFinancialData(val, filingType10K)
		filing.Date = key
		company.Reports = append(company.Reports, filing)
		break
	}
	log.Println(company)
}
