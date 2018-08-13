package main

type filingType string
type filingDocType string
type finDataType string

var (
	filingType10Q filingType = "10-Q"
	filingType10K filingType = "10-K"

	filingDocOps filingDocType = "Operations"
	filingDocInc filingDocType = "Income"
	filingDocBS  filingDocType = "Assets"
	filingDocCF  filingDocType = "Cash Flow"
	filingDocEN  filingDocType = "Entity Info"
	filingDocIg  filingDocType = "Ignore"

	finDataSharesOutstanding finDataType = "Shares Outstanding"
)

func main() {

	filingLinks := getFilingLinks("AAPL", filingType10Q)
	for _, link := range filingLinks {
		getFilingPage(link, filingType10Q)
		break
	}

}
