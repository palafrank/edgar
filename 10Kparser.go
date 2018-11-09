package edgar

import (
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func map10KReports(page io.Reader, filingLinks []string) map[filingDocType]string {
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
					docType := lookupDocType(token.String())
					if docType != filingDocIg {
						//Get the report number
						//fmt.Println("Found a wanted doc ", docType, token.String(), reportNum)
						_, ok := retData[docType]
						if !ok {
							retData[docType] = filingLinks[reportNum-1]
						}
						if len(retData) == len(requiredDocTypes) {
							return retData
						}
					}
				} else if a.Key == "id" && a.Val == "menu_cat3" {
					//Gone too far. Menu category 3 is beyond consolidated statements.
					//Stop parsing
					break loop
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
