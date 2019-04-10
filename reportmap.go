package edgar

import (
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type menuCat string

var (
	menuCatCover   menuCat = "Cover"
	menuCatFS      menuCat = "Financial statements"
	menuCatNFS     menuCat = "Notes to Financial statements"
	menuCatUnknown menuCat = "Unknown"
)

func getMenuCategory(data string) menuCat {
	data = strings.ToLower(data)
	if strings.Contains(data, "financial") && strings.Contains(data, "statement") {
		if strings.Contains(data, "note") {
			return menuCatNFS
		}
		return menuCatFS
	} else if strings.Contains(data, "cover") {
		return menuCatCover
	}
	return menuCatUnknown
}

func lookupDocType(data string, menu menuCat) filingDocType {

	data = strings.ToUpper(data)

	if menu == menuCatCover && strings.Contains(data, "DOCUMENT") &&
		strings.Contains(data, "ENTITY") {
		//Entity document
		return filingDocEN
	} else if menu == menuCatFS {
		if strings.Contains(data, "PARENTHETICAL") {
			//skip this doc
			return filingDocIg
		}
		// Financial statements
		if strings.Contains(data, "BALANCE SHEET") {
			//Balance sheet
			return filingDocBS
		} else if strings.Contains(data, "FINANCIAL POSITION") {
			//Balance sheet
			return filingDocBS
		} else if strings.Contains(data, "OPERATIONS") {
			//Operations statement
			return filingDocOps
		} else if strings.Contains(data, "INCOME") {
			//Income statement
			return filingDocInc
		} else if strings.Contains(data, "EARNINGS") {
			//Income statement
			return filingDocInc
		} else if strings.Contains(data, "CASH FLOWS") {
			//Cash flow statement
			return filingDocCF
		}
	} else if menu == menuCatNFS {
		// Notes to Financial statements
		if strings.Contains(data, "EARNINGS") && strings.Contains(data, "SHARE") {
			return filingDocEPSNotes
		} else if strings.Contains(data, "SHAREHOLDER") && strings.Contains(data, "EQUITY") {
			return filingDocEquity
		} else if strings.Contains(data, "DEBT") {
			return filingDocDebt
		}
	}
	return filingDocIg
}

func getMissingDocs(data map[filingDocType]string) string {

	if len(data) >= len(requiredDocTypes) {
		return ""
	}
	var diff []filingDocType
	for key := range requiredDocTypes {
		if _, ok := data[key]; !ok {
			switch key {
			case filingDocOps:
				if _, ok := data[filingDocInc]; ok {
					continue
				}
			case filingDocInc:
				if _, ok := data[filingDocOps]; ok {
					continue
				}
			}
			diff = append(diff, key)
		}
	}
	if len(diff) == 0 {
		return ""
	}

	var ret string
	ret = "[ "
	for _, val := range diff {
		ret = ret + " " + string(val)
	}
	ret += " ]"
	return ret
}

func mapReports(page io.Reader, filingLinks []string) map[filingDocType]string {

	menuCategory := menuCatUnknown

	retData := make(map[filingDocType]string)

	z := html.NewTokenizer(page)
	tt := z.Next()
loop:
	for tt != html.ErrorToken {
		token := z.Token()
		if token.Data == "a" {
			for _, a := range token.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "loadReport") {
					strs := strings.Split(a.Val, "loadReport")
					strs[1] = strings.Trim(strs[1], ";")
					reportNum, _ := strconv.Atoi(strings.Trim(strs[1], "()"))
					tt = z.Next() //Contains the text that describes the report
					if tt != html.TextToken {
						break
					}
					token = z.Token()
					docType := lookupDocType(token.String(), menuCategory)
					if docType != filingDocIg {
						//Get the report number
						_, ok := retData[docType]
						if !ok {
							retData[docType] = filingLinks[reportNum-1]
						}
					}
				} else if a.Key == "id" && strings.Contains(a.Val, "menu_cat") {
					// Set the menu level
					for !(token.Data == "a" && token.Type == html.EndTagToken) {
						if token.Type == html.TextToken {
							str := strings.TrimSpace(token.String())
							menuCategory = getMenuCategory(str)
						}
						z.Next()
						token = z.Token()
					}
					if menuCategory == menuCatUnknown {
						//Gone too far. Menu category 4 is beyond notes of financial statements.
						//Stop parsing
						break loop
					}
				}
			}
		}
		tt = z.Next()
	}
	ret := getMissingDocs(retData)
	if ret != "" {
		log.Println("Did not find the following filing documents: " + ret)
	}
	return retData
}
