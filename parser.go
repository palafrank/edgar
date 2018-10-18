package edgar

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var filingButtonId string = "interactiveDataBtn"

/*
  This is the parsing of query page where we get the list of filings of a given types
  ex: https://www.sec.gov/cgi-bin/browse-edgar?CIK=AAPL&owner=exclude&action=getcompany&type=10-Q&count=1&dateb=
  Assumptions of the parser:
  - There is interactive data available and there is a button that allows the user to click it
  - Since it is a link the tag will be a hyperlink with a button with the id=interactiveDataBtn
  - The actual link is the href attribute in the "a" token just before the id attribute
*/
func queryPageParser(page io.Reader, docType FilingType) map[string]string {

	filingInfo := make(map[string]string)

	z := html.NewTokenizer(page)

	data, err := parseTableRow(z, true)
	for err == nil {
		//This check for filing type will drop AMEND filings
		if len(data) == 5 && data[0] == string(docType) {
			//Drop filings before 2010
			year := getYear(data[3])
			if year >= thresholdYear {
				filingInfo[data[3]] = data[1]
			}
		}
		data, err = parseTableRow(z, true)
	}
	return filingInfo
}

/*
  The filing page parser
  - The top of the page has a list of reports.
  - Get all the reports (link to all the reports) and put it in an array
  - The Accordian on the side of the page identifies what each report is
  - Get the text of the accordian and map the type of the report to the report
  - Create a map of the report to report link
*/
func filingPageParser(page io.Reader, fileType FilingType) map[filingDocType]string {
	var filingLinks []string
	r := bufio.NewReader(page)
	s, e := r.ReadString('\n')

	for e == nil {
		//Get the number of reports available
		if strings.Contains(s, "var reports") == true {
			s1 := strings.Split(s, "(")
			s2 := strings.Split(s1[1], ")")
			cnt, _ := strconv.Atoi(s2[0])

			//cnt-1 because we skip the 'all' in the list
			for i := 0; i < cnt-1; i++ {
				s, e = r.ReadString('\n')
				s1 := strings.Split(s, " = ")
				s2 := strings.Split(s1[1], ";")
				s3 := strings.Trim(s2[0], "\"")
				s4 := strings.Split(s3, ".")
				s5 := s3
				//Sometimes the report is listed as an xml file??
				if s4[1] == "xml" {
					s5 = s4[0] + ".htm"
				}
				if !strings.Contains(s5, "htm") {
					panic("Dont know this type of report")
				}
				filingLinks = append(filingLinks, s5)
			}

			break
		}
		s, e = r.ReadString('\n')

	}

	switch fileType {
	case FilingType10K:
		log.Println("Getting 10K filing documents: ")
		docs := map10KReports(page, filingLinks)
		return docs
	case FilingType10Q:
		log.Println("Getting 10Q filing documents")
		docs := map10QReports(page, filingLinks)
		return docs
	}

	return nil

}

func parseTableData(z *html.Tokenizer, parseHref bool) string {
	token := z.Token()

	if token.Type != html.StartTagToken && token.Data != "td" {
		log.Fatal("Tokenizer passed incorrectly to parseTableData")
		return ""
	}

	for !(token.Data == "td" && token.Type == html.EndTagToken) {
		if token.Type == html.ErrorToken {
			break
		}

		if parseHref && token.Data == "a" && token.Type == html.StartTagToken {
			str := parseHyperLinkTag(z, token)
			if len(str) > 0 {
				return str
			}
		} else {
			if token.Type == html.TextToken {
				str := strings.TrimSpace(token.String())
				if len(str) > 0 {
					return str
				}
			}
		}
		//Going for the end of the td tag
		z.Next()
		token = z.Token()
	}
	return ""
}

func parseTableRow(z *html.Tokenizer, parseHref bool) ([]string, error) {
	var retData []string
	//Get the current token
	token := z.Token()

	//Check if this is really a table row
	for !(token.Type == html.StartTagToken && token.Data == "tr") {
		tt := z.Next()
		if tt == html.ErrorToken {
			return nil, errors.New("Done with parsing")
		}
		token = z.Token()
	}
	//Till the end of the row collect data from each data block
	for !(token.Data == "tr" && token.Type == html.EndTagToken) {

		if token.Type == html.ErrorToken {
			return nil, errors.New("Done with parsing")
		}
		if token.Data == "td" && token.Type == html.StartTagToken {
			parseFlag := parseHref
			//If the data is a number class just get the text = number
			for _, a := range token.Attr {
				if a.Key == "class" && (a.Val == "nump" || a.Val == "num") {
					parseFlag = false
				}
			}
			str := parseTableData(z, parseFlag)
			if len(str) > 0 {
				retData = append(retData, str)
			}
		}
		z.Next()
		token = z.Token()
	}

	return retData, nil
}

var reqHyperLinks = map[string]bool{
	"interactiveDataBtn": true,
}

func parseHyperLinkTag(z *html.Tokenizer, token html.Token) string {
	var href string
	var onclick string
	var id string

	for _, a := range token.Attr {
		switch a.Key {
		case "id":
			id = a.Val
		case "href":
			href = a.Val
		case "onclick":
			onclick = a.Val
			if str, err := getFinDataXBRLTag(onclick); err == nil {
				return str
			}
		}
	}

	text := ""
	//Finish up the hyperlink
	for !(token.Data == "a" && token.Type == html.EndTagToken) {
		/*
			if token.Type == html.TextToken {
				str := strings.TrimSpace(token.String())
				if len(str) > 0 {
					text = str
				}
			}
		*/
		z.Next()
		token = z.Token()
	}

	if _, ok := reqHyperLinks[id]; ok {
		return href
	}

	return text
}

func parseTableTitle(z *html.Tokenizer) []string {

	var strs []string
	token := z.Token()

	if token.Type != html.StartTagToken && token.Data != "th" {
		log.Fatal("Tokenizer passed incorrectly to parseTableData")
		return strs
	}

	for !(token.Data == "th" && token.Type == html.EndTagToken) {
		if token.Type == html.ErrorToken {
			break
		}

		if token.Type == html.TextToken {
			str := strings.TrimSpace(token.String())
			if len(str) > 0 {
				strs = append(strs, str)
			}
		}
		//Going for the end of the td tag
		z.Next()
		token = z.Token()
	}
	return strs
}

func parseTableHeading(z *html.Tokenizer) ([]string, error) {
	var retData []string
	//Get the current token
	token := z.Token()

	//Check if this is really a table row
	for !(token.Type == html.StartTagToken && token.Data == "tr") {
		tt := z.Next()
		if tt == html.ErrorToken {
			return nil, errors.New("Done with parsing")
		}
		token = z.Token()
	}

	//Till the end of the row collect data from each data block
	for !(token.Data == "tr" && token.Type == html.EndTagToken) {

		if token.Type == html.ErrorToken {
			return nil, errors.New("Done with parsing")
		}
		if token.Data == "th" && token.Type == html.StartTagToken {
			str := parseTableTitle(z)
			if len(str) > 0 {
				retData = append(retData, str...)
			}
		}
		z.Next()
		token = z.Token()
	}

	return retData, nil
}

func parseFilingScale(z *html.Tokenizer) map[scaleEntity]scaleFactor {
	scales := make(map[scaleEntity]scaleFactor)
	data, err := parseTableHeading(z)
	if err == nil {
		if len(data) > 0 {
			scales = filingScale(data)
		}
	}
	return scales
}
