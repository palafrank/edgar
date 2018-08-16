package main

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
func queryPageParser(page io.Reader, docType filingType) map[string]string {

	filingInfo := make(map[string]string)

	z := html.NewTokenizer(page)

	data, err := parseTableRow(z, true)
	for err == nil {
		if len(data) == 5 && data[0] == string(docType) {
			filingInfo[data[3]] = data[1]
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
func filingPageParser(page io.Reader, fileType filingType) map[filingDocType]string {
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
				filingLinks = append(filingLinks, s3)
			}

			break
		}
		s, e = r.ReadString('\n')

	}

	switch fileType {
	case filingType10K:
	case filingType10Q:
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

		if token.Type == html.TextToken {
			str := strings.TrimSpace(token.String())
			if len(str) > 0 {
				return str
			}
		}
		if parseHref {
			if token.Data == "a" && token.Type == html.StartTagToken {
				str := parseHyperLinkTag(z, token)
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
			str := parseTableData(z, parseHref)
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
	var id string
	var text string

	for _, a := range token.Attr {
		switch a.Key {
		case "id":
			id = a.Val
		case "href":
			href = a.Val
		}
	}
	for !(token.Data == "a" && token.Type == html.EndTagToken) {
		if token.Type == html.ErrorToken {
			log.Fatal("Unexpected table data parsing error")
			return ""
		}
		if token.Type == html.TextToken {
			text = token.Data
		}
		z.Next()
		token = z.Token()
	}
	if len(id) == 0 && len(text) == 0 {
		log.Fatal("Bad parsing of hyperlink tag")
		return ""
	} else {
		if len(id) > 0 {
			text = id
		}
	}
	_, ok := reqHyperLinks[text]
	if ok {
		return href
	}
	return ""
}
