package edgar

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var sampleTableRow = `<tr><td nowrap="nowrap">10-Q</td><td nowrap="nowrap"><a href="/Archives/edgar/data/320193/000032019318000100/0000320193-18-000100-index.htm" id="documentsbutton">&nbsp;Documents</a>&nbsp; <a href="/cgi-bin/viewer?action=view&amp;cik=320193&amp;accession_number=0000320193-18-000100&amp;xbrl_type=v" id="interactiveDataBtn">&nbsp;Interactive Data</a></td><td class="small" >Quarterly report [Sections 13 or 15(d)]<br />Acc-no: 0000320193-18-000100&nbsp;(34 Act)&nbsp; Size: 9 MB            </td><td>2018-08-01</td><td nowrap="nowrap"><a href="/cgi-bin/browse-edgar?action=getcompany&amp;filenum=001-36743&amp;owner=exclude&amp;count=10">001-36743</a><br>18985212         </td></tr><tr class="blueRow">`

var sampleRowWithXBRL = `<tr class="reu"><td class="pl " style="border-bottom: 0px;" valign="top"><a class="a" href="javascript:void(0);" onclick="top.Show.showAR( this, 'defref_us-gaap_StockholdersEquity', window );">Total shareholders&#8217; equity</a></td><td class="nump">134,047<span></span>
</td><td class="nump">128,249<span></span></td></tr>`

var sampleRowWithNumInLink = `<tr class="re">
        <td class="pl " style="border-bottom: 0px;" valign="top"><a class="a" href="javascript:void(0);" onclick="top.Show.showAR( this, 'defref_dei_EntityCommonStockSharesOutstanding', window );">Entity Common Stock, Shares Outstanding</a></td>
        <td class="text">&#xA0;<span></span></td>
        <td class="nump"><a title="dei_EntityCommonStockSharesOutstanding" onclick="toggleNextSibling(this);">266,252,295</a><span style="display:none;white-space:normal;text-align:left;">dei_EntityCommonStockSharesOutstanding</span><span></span></td>
        <td class="text">&#xA0;<span></span></td>
      </tr>
`

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

func TestParsingXBRLDef(t *testing.T) {
	page := strings.NewReader(sampleRowWithXBRL)
	z := html.NewTokenizer(page)
	data, err := parseTableRow(z, true)
	if err != nil {
		t.Error("Parser returned error while parsing XBRL: " + err.Error())
		return
	}
	if len(data) != 3 {
		t.Error("Parser returned unexpected number of data items: " + string(len(data)))
		return
	}
	if data[0] != "defref_us-gaap_StockholdersEquity" {
		t.Error("Did not get the expected financial data tag: ", data[0])
	}
	if data[1] != "134,047" {
		t.Error("Did not get the righ value from the table: ", data[1])
	}
	if data[2] != "128,249" {
		t.Error("Did not get the righ value from the table: ", data[2])
	}
}

func TestParsingNumInLink(t *testing.T) {
	page := strings.NewReader(sampleRowWithNumInLink)
	z := html.NewTokenizer(page)
	data, err := parseTableRow(z, true)
	if err != nil {
		t.Error("Parser returned error while parsing XBRL: " + err.Error())
		return
	}
	if len(data) != 2 {
		t.Error("Parser returned unexpected number of data items: " + string(len(data)))
		return
	}
	if data[0] != "defref_dei_EntityCommonStockSharesOutstanding" {
		t.Error("Did not get the expected financial data tag: ", data[0])
	}
	if data[1] != "266,252,295" {
		t.Error("Did not get the righ value from the table: ", data[1])
	}
}

func TestGetCIK(t *testing.T) {
	cik := getCompanyCIK("MSFT")
	if cik != "0000789019" {
		t.Error("Incorrect CIK parser for MSFT - ", cik)
	}
	cik = getCompanyCIK("GE")
	if cik != "0000040545" {
		t.Error("Incorrect CIK parser for MSFT - ", cik)
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
	f, _ := os.Open("samples/sample_query.html")
	links := queryPageParser(f, FilingType10Q)
	f.Close()
	if len(links) != 10 {
		t.Error("Incorrect number of filing links found ", len(links))
	}

	for key, val := range links {
		if val != valid[key] {
			t.Error("Incorrect: ", key, val)
		}
	}

}

func TestParseCIKAndDocID(t *testing.T) {
	str1 := "/cgi-bin/viewer?action=view&cik=320193&accession_number=0001193125-15-259935&xbrl_type=v"
	s1, s2 := parseCikAndDocId(str1)
	if s1 != "320193" || s2 != "000119312515259935" {
		t.Error("Error in parsing CIK and doc id ", s1, s2)
	}
}

func TestFiling10QParser(t *testing.T) {
	var check = map[filingDocType]string{
		filingDocCF:  "/Archives/edgar/data/320193/000032019318000100/R7.htm",
		filingDocInc: "/Archives/edgar/data/320193/000032019318000100/R3.htm",
		filingDocEN:  "/Archives/edgar/data/320193/000032019318000100/R1.htm",
		filingDocOps: "/Archives/edgar/data/320193/000032019318000100/R2.htm",
		filingDocBS:  "/Archives/edgar/data/320193/000032019318000100/R5.htm",
	}
	f, _ := os.Open("samples/sample_10Q.html")
	docs := filingPageParser(f, FilingType10Q)
	f.Close()
	for key, val := range check {
		if docs[key] != val {
			t.Error("Did not get the expected number of filing document in the 10K")
		}
	}
}

func TestFiling10KParser(t *testing.T) {
	var check = map[filingDocType]string{
		filingDocCF:  "/Archives/edgar/data/320193/000119312515356351/R8.htm",
		filingDocInc: "/Archives/edgar/data/320193/000119312515356351/R3.htm",
		filingDocEN:  "/Archives/edgar/data/320193/000119312515356351/R1.htm",
		filingDocOps: "/Archives/edgar/data/320193/000119312515356351/R2.htm",
		filingDocBS:  "/Archives/edgar/data/320193/000119312515356351/R5.htm",
	}
	f, _ := os.Open("samples/sample_10K.html")
	docs := filingPageParser(f, FilingType10K)
	f.Close()
	for key, val := range check {
		if docs[key] != val {
			t.Error("Did not get the expected number of filing document in the 10K")
		}
	}
}

func TestParsingReports(t *testing.T) {
	url := "cgi-bin/viewer?action=view&cik=789019&accession_number=0001193125-13-310206&xbrl_type=v"
	for i := 0; i < 1; i++ {
		report, err := getFinancialData(url, FilingType10K)
		if err != nil {
			t.Error("Failed to parse financial data: ", err.Error())
			return
		}
		if report.Entity.ShareCount != 8329956402 {
			t.Error("Incorrect sharcount parsed ", report.Entity.ShareCount)
		}
		if report.Ops.Revenue != 77849000000 {
			t.Error("Incorrect revenue parsed ", report.Ops.Revenue)
		}
		if report.Ops.CostOfSales != 20249000000 {
			t.Error("Incorrect cost of sales parsed ", report.Ops.CostOfSales)
		}
		if report.Ops.GrossMargin != 57600000000 {
			t.Error("Incorrect gross margin parsed ", report.Ops.GrossMargin)
		}
		if report.Ops.OpIncome != 26764000000 {
			t.Error("Incorrect ops income parsed ", report.Ops.OpIncome)
		}
		if report.Ops.OpExpense != 30836000000 {
			t.Error("Incorrect ops expense parsed ", report.Ops.OpExpense)
		}
		if report.Ops.NetIncome != 21863000000 {
			t.Error("Incorrect net income parsed ", report.Ops.NetIncome)
		}
		if report.Cf.OpCashFlow != 28833000000 {
			t.Error("Incorrect operating cashflow parsed ", report.Cf.OpCashFlow)
		}
		if report.Cf.CapEx != -4257000000 {
			t.Error("Incorrect capex parsed ", report.Cf.CapEx)
		}
		if report.Bs.LDebt != 12601000000 {
			t.Error("Incorrect long term debt parsed ", report.Bs.LDebt)
		}
		if report.Bs.SDebt != 0 {
			t.Error("Incorrect short term debt parsed ", report.Bs.SDebt)
		}
		if report.Bs.CLiab != 37417000000 {
			t.Error("Incorrect current liabilities parsed ", report.Bs.CLiab)
		}
		if report.Bs.Deferred != 20639000000 {
			t.Error("Incorrect deferred revenue parsed ", report.Bs.Deferred)
		}
		if report.Bs.Retained != 9895000000 {
			t.Error("Incorrect retained earnings parsed ", report.Bs.Retained)
		}
		if report.Bs.Equity != 78944000000 {
			t.Error("Incorrect shareholder equity parsed ", report.Bs.Equity)
		}
	}
}

func TestFiling10KParser1(t *testing.T) {
	var check = map[filingDocType]string{
		filingDocCF:  "/Archives/edgar/data/320193/000119312511282113/R6.htm",
		filingDocEN:  "/Archives/edgar/data/320193/000119312511282113/R1.htm",
		filingDocOps: "/Archives/edgar/data/320193/000119312511282113/R2.htm",
		filingDocBS:  "/Archives/edgar/data/320193/000119312511282113/R3.htm",
	}
	f, _ := os.Open("samples/sample_10K_1.html")
	docs := filingPageParser(f, FilingType10K)
	f.Close()
	for key, val := range check {
		if docs[key] != val {
			t.Error("Did not get the expected number of filing document in the 10K")
		}
	}

}

func TestEntityParser(t *testing.T) {
	f, _ := os.Open("samples/sample_entity.html")
	entity := new(entityData)
	_, err := reportParser(f, entity)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else if entity.ShareCount != 4829926000 {
		t.Error("Incorrect sharecount value parsed")
	}
}

func TestEntity1Parser(t *testing.T) {
	f, _ := os.Open("samples/sample_entity1.html")
	entity := new(entityData)
	_, err := reportParser(f, entity)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else if entity.ShareCount != 266252295 {
		t.Error("Incorrect sharecount value parsed: ", entity.ShareCount)
	}
}

func Test10KEntityParser(t *testing.T) {
	f, _ := os.Open("samples/sample_10K_entity.html")
	entity := new(entityData)
	_, err := reportParser(f, entity)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else if entity.ShareCount != 5575331000 {
		t.Error("Incorrect sharecount value parsed")
	}
}

func TestOpsParser(t *testing.T) {
	f, _ := os.Open("samples/sample_ops.html")
	ops := new(opsData)
	_, err := reportParser(f, ops)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if ops.Revenue != 53265000000 {
			t.Error("Revenue amount did not match")
		}
		if ops.CostOfSales != 32844000000 {
			t.Error("Cost of Sales amount did not match")
		}
		if ops.GrossMargin != 20421000000 {
			t.Error("Gross margin amount did not match ", ops.GrossMargin)
		}
		if ops.OpExpense != 7809000000 {
			t.Error("Operational Expense amount did not match")
		}
		if ops.OpIncome != 12612000000 {
			t.Error("Operational Income amount did not match")
		}
		if ops.NetIncome != 11519000000 {
			t.Error("Net income amount did not match " + strconv.Itoa(int(ops.NetIncome)))
		}
	}
}

func Test10KOpsParser(t *testing.T) {
	f, _ := os.Open("samples/sample_10K_ops.html")
	ops := new(opsData)
	_, err := reportParser(f, ops)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if ops.Revenue != 233715000000 {
			t.Error("Revenue amount did not match")
		}
		if ops.CostOfSales != 140089000000 {
			t.Error("Cost of Sales amount did not match")
		}
		if ops.GrossMargin != 93626000000 {
			t.Error("Gross margin amount did not match")
		}
		if ops.OpExpense != 22396000000 {
			t.Error("Operational Expense amount did not match")
		}
		if ops.OpIncome != 71230000000 {
			t.Error("Operational Income amount did not match")
		}
		if ops.NetIncome != 53394000000 {
			t.Error("Net income amount did not match")
		}
	}
}

func TestCfParser(t *testing.T) {
	f, _ := os.Open("samples/sample_cf.html")
	cf := new(cfData)
	_, err := reportParser(f, cf)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if cf.OpCashFlow != 57911000000 {
			t.Error("Incorrect cash flow from operations value parsed")
		}
		if cf.CapEx != float64(-10272000000) {
			t.Error("Incorrect capital expenditure value parsed")
		}
	}
}

func Test10KCfParser(t *testing.T) {
	f, _ := os.Open("samples/sample_10K_cf.html")
	cf := new(cfData)
	_, err := reportParser(f, cf)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if cf.OpCashFlow != 81266000000 {
			t.Error("Incorrect cash flow from operations value parsed")
		}
		if cf.CapEx != float64(-11247000000) {
			t.Error("Incorrect capital expenditure value parsed ", cf.CapEx)
		}
	}
}

func TestBSParser(t *testing.T) {
	f, _ := os.Open("samples/sample_bs.html")
	bs := new(bsData)
	_, err := reportParser(f, bs)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if bs.CLiab != 88548000000 {
			t.Error("Incorrect current liabilities from balance sheet value parsed")
		}
		if bs.LDebt != 97128000000 {
			t.Error("Incorrect long term debt from balance sheet value parsed")
		}
		/*
			if bs.SDebt != 5498000000 {
				t.Error("Incorrect short term debt from balance sheet value parsed")
			}
		*/
		if bs.Retained != 79436000000 {
			t.Error("Incorrect retained earningd from balance sheet value parsed")
		}
	}
}

func TestBS1Parser(t *testing.T) {
	f, _ := os.Open("samples/sample_bs1.html")
	bs := new(bsData)
	_, err := reportParser(f, bs)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if bs.CLiab != 5018000000 {
			t.Error("Incorrect current liabilities: ", bs.CLiab)
		}
		if bs.LDebt != 14846000000 {
			t.Error("Incorrect long term debt ", bs.LDebt)
		}

		if bs.Deferred != 27000000 {
			t.Error("Incorrect Deferred ", bs.Deferred)
		}

		if bs.Equity != 28331000000 {
			t.Error("Incorrect equity ", bs.Equity)
		}

		if bs.Retained != -198000000 {
			t.Error("Incorrect retained earningd ", bs.Retained)
		}
	}
}

func Test10KBSParser(t *testing.T) {
	f, _ := os.Open("samples/sample_10K_bs.html")
	bs := new(bsData)
	_, err := reportParser(f, bs)
	f.Close()
	if err != nil {
		t.Error(err.Error())
	} else {
		if bs.CLiab != 80610000000 {
			t.Error("Incorrect current liabilities from balance sheet value parsed")
		}
		if bs.LDebt != 53463000000 {
			t.Error("Incorrect long term debt from balance sheet value parsed")
		}
		/*
			if bs.SDebt != 2500000000 {
				t.Error("Incorrect short term debt from balance sheet value parsed")
			}
		*/
		if bs.Retained != 92284000000 {
			t.Error("Incorrect retained earningd from balance sheet value parsed")
		}
	}
}

func TestFinReportMarshal(t *testing.T) {

	var file filing
	var data financialReport

	file.Date = "2017-02-1"
	file.Company = "AAPL"
	file.FinData = &data

	comp := newCompany("AAPL")
	comp.AddReport(&file)

	data.DocType = FilingType10K
	f, _ := os.Open("samples/sample_10K_bs.html")
	data.Bs = new(bsData)
	_, _ = reportParser(f, data.Bs)
	f.Close()
	f, _ = os.Open("samples/sample_10K_cf.html")
	data.Cf = new(cfData)
	_, _ = reportParser(f, data.Cf)
	f.Close()
	f, _ = os.Open("samples/sample_10K_ops.html")
	data.Ops = new(opsData)
	_, _ = reportParser(f, data.Ops)
	f.Close()
	f, _ = os.Open("samples/sample_10K_entity.html")
	data.Entity = new(entityData)
	_, _ = reportParser(f, data.Entity)
	f.Close()
	str := data.String()
	if !(strings.Contains(str, "Entity Information") &&
		strings.Contains(str, "Operational Information") &&
		strings.Contains(str, "Balance Sheet Information") &&
		strings.Contains(str, "Cash Flow Information")) {
		t.Error("Error generating the JSON document for financial report")
	}
	f, _ = os.Open("samples/sample_10K_marshal.json")
	b, _ := ioutil.ReadAll(f)
	f.Close()

	//There is an extra byte at the end of the save file that needs to be
	//eliminated to avoid a mismatch
	if str != string(b[:len(b)-1]) {
		t.Error("Marshaled data doesnot match reference JSON\n", str)
	}
}

func TestFolderReader(t *testing.T) {
	f, _ := os.Open("samples/sample_folder.json")
	fetcher := NewFilingFetcher()
	c, err := fetcher.CreateFolder(f)
	if err != nil {
		t.Error(err)
		return
	}
	f.Close()
	f, _ = os.Open("samples/sample_folder.json")
	b, _ := ioutil.ReadAll(f)
	f.Close()
	//There is an extra byte at the end of the save file that needs to be
	//eliminated to avoid a mismatch
	if c.String() != string(b[:len(b)-1]) {
		t.Error("Created folder does not match sample stored folder\n", c.String())
	}
}

// LIVE TESTS:
//     These tests are run live against EDGAR website. They are commented out
//     to avoid hitting the site during repeated unit testing.
//     Uncomment them when a live test is needed to verify something that is
//     not covered in the samples.

func TestFolderWriter(t *testing.T) {
	fetcher := NewFilingFetcher()
	c, err := fetcher.CompanyFolder("AGN", FilingType10K)
	if err != nil {
		t.Error(err)
	}
	files := c.AvailableFilings(FilingType10K)
	for _, val := range files {
		if val.Year() == 2018 || val.Year() == 2017 {
			c.Filing(FilingType10K, val)
		}
	}
	f, _ := os.Open("samples/sample_writer.json")
	b, _ := ioutil.ReadAll(f)
	f.Close()
	//There is an extra byte at the end of the save file that needs to be
	//eliminated to avoid a mismatch
	if c.String() != string(b[:len(b)-1]) {
		t.Error("Created folder does not match sample stored folder ", c.String())
	}
}

func TestLiveMSFTParsing(t *testing.T) {
	fmt.Println("*** Running a live test ***")
	fetcher := NewFilingFetcher()
	c, err := fetcher.CompanyFolder("MSFT", FilingType10K)
	if err != nil {
		t.Error(err)
	}
	files := c.AvailableFilings(FilingType10K)
	for _, val := range files {
		if val.Year() == 2018 || val.Year() == 2015 || val.Year() == 2011 {
			fs, err := c.Filing(FilingType10K, val)
			if err != nil {
				t.Error("Failed to get filing " + val.String())
			}
			if val.Year() == 2018 {
				if val, _ := fs.LongTermDebt(); val != 72242000000 {
					t.Error("Incorrect Long term debt-2018: ", val)
				}
				if val, _ := fs.RetainedEarnings(); val != 13682000000 {
					t.Error("Incorrect retained earnings-2018: ", val)
				}
				if val, _ := fs.ShareCount(); val != 7668217316 {
					t.Error("Incorrect share count-2018: ", val)
				}
				if val, _ := fs.OperatingCashFlow(); val != 43884000000 {
					t.Error("Incorrect share count-2018: ", val)
				}
			}
			if val.Year() == 2011 {
				if val, _ := fs.LongTermDebt(); val != 11921000000 {
					t.Error("Incorrect Long term debt-2018: ", val)
				}
				if val, _ := fs.RetainedEarnings(); val != -6332000000 {
					t.Error("Incorrect retained earnings-2018: ", val)
				}
				if val, _ := fs.ShareCount(); val != 8378265782 {
					t.Error("Incorrect share count-2018: ", val)
				}
				if val, _ := fs.OperatingCashFlow(); val != 26994000000 {
					t.Error("Incorrect share count-2018: ", val)
				}
				if val, _ := fs.CapitalExpenditure(); val != -2355000000 {
					t.Error("Incorrect share count-2018: ", val)
				}
			}
			if val.Year() == 2015 {
				if val, _ := fs.LongTermDebt(); val != 27808000000 {
					t.Error("Incorrect Long term debt-2018: ", val)
				}
				if val, _ := fs.RetainedEarnings(); val != 9096000000 {
					t.Error("Incorrect retained earnings-2018: ", val)
				}
				if val, _ := fs.ShareCount(); val != 7997980969 {
					t.Error("Incorrect share count-2018: ", val)
				}
				if val, _ := fs.OperatingCashFlow(); val != 29080000000 {
					t.Error("Incorrect share count-2018: ", val)
				}
				if val, _ := fs.CapitalExpenditure(); val != -5944000000 {
					t.Error("Incorrect share count-2018: ", val)
				}
			}
		}
	}
}
