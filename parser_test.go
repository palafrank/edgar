package main

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var sampleTableRow = `<tr><td nowrap="nowrap">10-Q</td><td nowrap="nowrap"><a href="/Archives/edgar/data/320193/000032019318000100/0000320193-18-000100-index.htm" id="documentsbutton">&nbsp;Documents</a>&nbsp; <a href="/cgi-bin/viewer?action=view&amp;cik=320193&amp;accession_number=0000320193-18-000100&amp;xbrl_type=v" id="interactiveDataBtn">&nbsp;Interactive Data</a></td><td class="small" >Quarterly report [Sections 13 or 15(d)]<br />Acc-no: 0000320193-18-000100&nbsp;(34 Act)&nbsp; Size: 9 MB            </td><td>2018-08-01</td><td nowrap="nowrap"><a href="/cgi-bin/browse-edgar?action=getcompany&amp;filenum=001-36743&amp;owner=exclude&amp;count=10">001-36743</a><br>18985212         </td></tr><tr class="blueRow">`

func TestParsingTableRow(t *testing.T) {
	page := strings.NewReader(sampleTableRow)
	z := html.NewTokenizer(page)
	data, err := parseTableRow(z, true)
	if err != nil {
		t.Error("Error parsing the table row with href enabled")
	}
	if len(data) != 5 {
		t.Error("Incorrect number of columns parsed in table row parsing")
	}
	if data[1] != "/cgi-bin/viewer?action=view&cik=320193&accession_number=0000320193-18-000100&xbrl_type=v" {
		t.Error("Incorrect parsing of HREF in the table row")
	}
	if data[0] != "10-Q" {
		t.Error("Incorrect parsing of the document type")
	}
	if data[3] != "2018-08-01" {
		t.Error("Incorrect date extracted while parsing table row")
	}
}

func TestFilingQuery(t *testing.T) {
	valid := map[string]string{
		"2018-08-01": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0000320193-18-000100&xbrl_type=v",
		"2018-05-02": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0000320193-18-000070&xbrl_type=v",
		"2018-02-02": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0000320193-18-000007&xbrl_type=v",
		"2017-08-02": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0000320193-17-000009&xbrl_type=v",
		"2017-05-03": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001628280-17-004790&xbrl_type=v",
		"2017-02-01": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001628280-17-000717&xbrl_type=v",
		"2016-07-27": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001628280-16-017809&xbrl_type=v",
		"2016-04-27": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001193125-16-559625&xbrl_type=v",
		"2016-01-27": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001193125-16-439878&xbrl_type=v",
		"2015-07-22": "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001193125-15-259935&xbrl_type=v",
	}
	f, _ := os.Open("./sample_query.html")
	links := queryPageParser(f, filingType10Q)
	if len(links) != 10 {
		t.Error("Incorrect number of filing links found")
	}

	for key, val := range links {
		if val != valid[key] {
			t.Error("Incorrect link parsed from the query document")
		}
	}

}

func TestFilingParser(t *testing.T) {
	f, _ := os.Open("./sample_10Q.html")
	docs := filingPageParser(f, filingType10Q)
	if len(docs) != 5 {
		t.Error("Did not get the expected number of filing document in the 10Q")
	}
}

func TestEntityParser(t *testing.T) {
	f, _ := os.Open("./sample_entity.html")
	entity, err := getEntityData(f)
	if err != nil {
		t.Error(err.Error())
	} else if entity == nil {
		t.Error("Entity object was not returned")
	} else if entity.ShareCount != 4829926 {
		t.Error("Incorrect sharecount value parsed")
	}
}

func TestOpsParser(t *testing.T) {
	f, _ := os.Open("./sample_ops.html")
	ops, err := getOpsData(f)
	if err != nil {
		t.Error(err.Error())
	} else if ops == nil {
		t.Error("Operations object was not returned")
	} else {
		if ops.Revenue != 53265 {
			t.Error("Revenue amount did not match")
		}
		if ops.CostOfSales != 32844 {
			t.Error("Cost of Sales amount did not match")
		}
		if ops.GrossMargin != 20421 {
			t.Error("Gross margin amount did not match")
		}
		if ops.OpExpense != 7809 {
			t.Error("Operational Expense amount did not match")
		}
		if ops.OpIncome != 12612 {
			t.Error("Operational Income amount did not match")
		}
		if ops.NetIncome != 11519 {
			t.Error("Net income amount did not match")
		}
	}
}

func TestCfParser(t *testing.T) {
	f, _ := os.Open("./sample_cf.html")
	cf, err := getCfData(f)
	if err != nil {
		t.Error(err.Error())
	} else if cf == nil {
		t.Error("Cash flow object was not returned")
	} else {
		if cf.OpCashFlow != 57911 {
			t.Error("Incorrect cash flow from operations value parsed")
		}
		if cf.CapEx != int64(-10272) {
			t.Error("Incorrect capital expenditure value parsed")
		}
	}
}

func TestBSParser(t *testing.T) {
	f, _ := os.Open("./sample_bs.html")
	bs, err := getBSData(f)
	if err != nil {
		t.Error(err.Error())
	} else if bs == nil {
		t.Error("Balance sheet object was not returned")
	} else {
		if bs.CLiab != 88548 {
			t.Error("Incorrect current liabilities from balance sheet value parsed")
		}
		if bs.LDebt != 97128 {
			t.Error("Incorrect long term debt from balance sheet value parsed")
		}
		if bs.SDebt != 5498 {
			t.Error("Incorrect short term debt from balance sheet value parsed")
		}
		if bs.Retained != 79436 {
			t.Error("Incorrect retained earningd from balance sheet value parsed")
		}
	}
}
